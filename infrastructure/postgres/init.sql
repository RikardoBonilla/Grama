-- init.sql - Database initialization for Grama
-- Runs ONCE when the container volume is empty (docker-entrypoint-initdb.d).
-- Executed as POSTGRES_USER (grama_admin superuser).
-- Role and grant setup lives in 02_roles.sh (needs env vars).

-- Extensions ------------------------------------------------------------------
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";  -- uuid_generate_v4()
CREATE EXTENSION IF NOT EXISTS "pgcrypto";   -- gen_random_bytes(), crypt()
CREATE EXTENSION IF NOT EXISTS "btree_gist"; -- EXCLUDE constraint on reservations

-- Shared trigger function for updated_at --------------------------------------
-- Every table that has updated_at attaches this trigger.
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Audit log -------------------------------------------------------------------
-- Append-only table for sensitive operations (reservations, payments, role changes).
-- grama_user may INSERT but never UPDATE or DELETE records here.
CREATE TABLE IF NOT EXISTS audit_log (
    id         UUID        PRIMARY KEY DEFAULT uuid_generate_v4(),
    table_name TEXT        NOT NULL,
    record_id  UUID        NOT NULL,
    action     TEXT        NOT NULL CHECK (action IN ('INSERT', 'UPDATE', 'DELETE')),
    actor_id   UUID,
    old_data   JSONB,
    new_data   JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
