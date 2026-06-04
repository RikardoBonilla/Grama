#!/bin/bash
# init-roles.sh - Create the application role with least privilege.
# Runs after init.sql (02_roles.sh in docker-entrypoint-initdb.d, alphabetical order).
# Requires GRAMA_DB_PASSWORD to be set in the container environment.
set -euo pipefail

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
    -- Application role: no superuser, no createdb, no createrole.
    CREATE ROLE grama_user WITH LOGIN PASSWORD '${GRAMA_DB_PASSWORD}';

    -- Connect to the database and use the public schema.
    GRANT CONNECT ON DATABASE ${POSTGRES_DB} TO grama_user;
    GRANT USAGE ON SCHEMA public TO grama_user;

    -- Audit log: INSERT only - append-only by design.
    GRANT INSERT ON audit_log TO grama_user;

    -- Any table that grama_admin creates in the future auto-grants
    -- SELECT, INSERT, UPDATE to grama_user. No DELETE, no DDL.
    ALTER DEFAULT PRIVILEGES IN SCHEMA public
        GRANT SELECT, INSERT, UPDATE ON TABLES TO grama_user;
EOSQL
