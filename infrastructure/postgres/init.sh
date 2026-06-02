#!/usr/bin/env bash
# infrastructure/postgres/init.sh
#
# Executed ONCE by docker-entrypoint-initdb.d when the container is first created.
# Shell scripts in that directory have access to container environment variables;
# .sql files do NOT — that is why this is a .sh, not a .sql.
#
# NEVER add "set -x" here: it would print GRAMA_DB_PASSWORD in container logs.
set -euo pipefail

# Fail fast if required variables are missing.
: "${GRAMA_DB_USER:?GRAMA_DB_USER is required — set it in docker-compose environment}"
: "${GRAMA_DB_PASSWORD:?GRAMA_DB_PASSWORD is required — set it in docker-compose environment}"
: "${POSTGRES_DB:?POSTGRES_DB is required}"

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL

    -- ─────────────────────────────────────────────────────────────────────
    -- 1. EXTENSIONS  (require superuser)
    -- ─────────────────────────────────────────────────────────────────────
    -- uuid_generate_v4(): primary key for every table — UUIDs, never BIGSERIAL.
    CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

    -- pgcrypto: gen_random_bytes() for future server-side crypto helpers.
    -- Argon2id hashing is done in Go; this covers any DB-level crypto operations.
    CREATE EXTENSION IF NOT EXISTS "pgcrypto";

    -- btree_gist: enables GiST indexes on scalar types (UUID, TEXT) mixed with
    -- range types. Required by the anti-overbooking EXCLUDE constraint on
    -- the reservations table (Sprint 3 — GRAM-20).
    CREATE EXTENSION IF NOT EXISTS "btree_gist";

    -- ─────────────────────────────────────────────────────────────────────
    -- 2. APPLICATION ROLE  (minimum privilege)
    -- ─────────────────────────────────────────────────────────────────────
    -- grama_user is the ONLY role the Go backend uses at runtime.
    -- Allowed: SELECT, INSERT, UPDATE.
    -- Denied:  DELETE, DDL (DROP, ALTER, CREATE), TRUNCATE.
    --
    -- Why no DELETE: if the app is compromised, the attacker cannot bulk-erase
    -- reservations or payment records. Audit trail stays intact.
    -- Soft deletes and status changes (e.g. cancelled) use UPDATE instead.
    DO \$\$
    BEGIN
        IF NOT EXISTS (SELECT FROM pg_roles WHERE rolname = '${GRAMA_DB_USER}') THEN
            CREATE ROLE "${GRAMA_DB_USER}" WITH LOGIN PASSWORD '${GRAMA_DB_PASSWORD}';
        ELSE
            -- Rotate password if .env changed and the container was re-created.
            ALTER ROLE "${GRAMA_DB_USER}" WITH PASSWORD '${GRAMA_DB_PASSWORD}';
        END IF;
    END
    \$\$;

    GRANT CONNECT ON DATABASE "${POSTGRES_DB}" TO "${GRAMA_DB_USER}";
    GRANT USAGE ON SCHEMA public TO "${GRAMA_DB_USER}";

    -- Apply to ALL tables and sequences created by grama_admin in the future
    -- (migrations). Does NOT cover tables created earlier in this same session —
    -- audit_log receives an explicit GRANT below.
    ALTER DEFAULT PRIVILEGES IN SCHEMA public
        GRANT SELECT, INSERT, UPDATE ON TABLES TO "${GRAMA_DB_USER}";

    ALTER DEFAULT PRIVILEGES IN SCHEMA public
        GRANT USAGE, SELECT ON SEQUENCES TO "${GRAMA_DB_USER}";

    -- ─────────────────────────────────────────────────────────────────────
    -- 3. REUSABLE TRIGGER FUNCTION
    -- ─────────────────────────────────────────────────────────────────────
    -- Every table with an updated_at column attaches this trigger.
    -- It fires BEFORE UPDATE so updated_at is always accurate, regardless
    -- of whether the application sets the field or not.
    -- Exception: audit_log is append-only and intentionally has NO updated_at.
    CREATE OR REPLACE FUNCTION update_updated_at_column()
    RETURNS TRIGGER
    LANGUAGE plpgsql
    AS \$func\$
    BEGIN
        NEW.updated_at = NOW();
        RETURN NEW;
    END;
    \$func\$;

    -- ─────────────────────────────────────────────────────────────────────
    -- 4. AUDIT LOG
    -- ─────────────────────────────────────────────────────────────────────
    -- Append-only record of sensitive business events.
    -- grama_user may only INSERT and SELECT — never UPDATE or DELETE.
    --
    -- Design decisions:
    --   - No updated_at / no trigger: this table is immutable by design.
    --   - action uses DOMAIN EVENT names (not SQL verbs) so the log is
    --     readable in compliance reports without needing a translation table.
    --   - actor_type distinguishes human users ('user') from automated
    --     system operations ('system', e.g. the reservation expiration job).
    --   - When actor_type = 'user', actor_id must be non-null: a missing
    --     actor_id on a user action is a bug, not a valid state.
    CREATE TABLE IF NOT EXISTS audit_log (
        id          UUID        NOT NULL DEFAULT uuid_generate_v4(),
        table_name  TEXT        NOT NULL,
        record_id   UUID        NOT NULL,
        action      TEXT        NOT NULL,
        actor_type  TEXT        NOT NULL DEFAULT 'user',
        actor_id    UUID,
        old_data    JSONB,
        new_data    JSONB,
        created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),

        CONSTRAINT audit_log_pkey PRIMARY KEY (id),

        -- Domain events. Extend this list when new sensitive operations are added.
        -- Deliberately NOT using SQL verbs (INSERT/UPDATE/DELETE) — those do not
        -- convey business meaning and make compliance queries harder to write.
        CONSTRAINT audit_log_action_chk CHECK (action IN (
            'USER_REGISTERED',
            'LOGIN_FAILED',
            'ROLE_CHANGED',
            'RESERVATION_CREATED',
            'RESERVATION_CANCELLED',
            'PRICE_CHANGED',
            'PAYMENT_REGISTERED',
            'PAYMENT_PROOF_UPLOADED',
            'PAYMENT_VERIFIED',
            'PAYMENT_REJECTED'
        )),

        CONSTRAINT audit_log_actor_type_chk CHECK (
            actor_type IN ('user', 'system')
        ),

        -- Human actions must have an identifiable actor.
        -- System actions (jobs, migrations) may omit actor_id.
        CONSTRAINT audit_log_actor_id_required CHECK (
            actor_type != 'user' OR actor_id IS NOT NULL
        )
    );

    COMMENT ON TABLE audit_log IS
        'Append-only record of sensitive business events (reservations, payments, '
        'role changes). grama_user: INSERT and SELECT only — never UPDATE or DELETE.';

    -- Explicit grant required: ALTER DEFAULT PRIVILEGES above does not cover
    -- tables created earlier in the same session.
    GRANT SELECT, INSERT ON audit_log TO "${GRAMA_DB_USER}";

    -- Query: "show all events for reservation X" (primary read pattern).
    CREATE INDEX IF NOT EXISTS idx_audit_log_record
        ON audit_log (table_name, record_id);

    -- Query: "what did user Y do today?" (compliance, incident response).
    -- Partial index skips system rows (actor_id IS NULL) to stay lean.
    CREATE INDEX IF NOT EXISTS idx_audit_log_actor
        ON audit_log (actor_id)
        WHERE actor_id IS NOT NULL;

EOSQL

echo "✓ grama_db initialized: extensions, grama_user role, trigger function, audit_log"
