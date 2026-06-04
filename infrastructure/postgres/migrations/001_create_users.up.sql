-- 001_create_users.up.sql
-- Sprint 1, GRAM-13: users table (Identity bounded context)
-- Run as grama_admin. Requires init.sql to have run (uuid-ossp, update_updated_at_column).

CREATE TABLE users (
    id               UUID        PRIMARY KEY DEFAULT uuid_generate_v4(),
    name             TEXT        NOT NULL,
    email            TEXT        NOT NULL,
    password_hash    TEXT        NOT NULL,
    role             TEXT        NOT NULL,
    consent_given    BOOLEAN     NOT NULL DEFAULT false,
    consent_given_at TIMESTAMPTZ,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT users_email_unique  UNIQUE (email),
    CONSTRAINT users_email_lower   CHECK (email = lower(email)),
    CONSTRAINT users_email_format  CHECK (email ~* '^[^@]+@[^@]+\.[^@]+$'),
    CONSTRAINT users_role_valid    CHECK (role IN ('owner', 'operator', 'client')),
    CONSTRAINT users_hash_nonempty CHECK (length(password_hash) > 0),
    CONSTRAINT users_name_nonempty CHECK (length(name) > 0)
);

-- Role index: authorization queries filter by role frequently.
CREATE INDEX idx_users_role       ON users (role);
-- Date index: occupancy and revenue reports group/filter by date.
CREATE INDEX idx_users_created_at ON users (created_at);

-- Keep updated_at current on every row update.
CREATE TRIGGER users_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Application role permissions: read, insert, update only. No DELETE, no DDL.
GRANT SELECT, INSERT, UPDATE ON users TO grama_user;
