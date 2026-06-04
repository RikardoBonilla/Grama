-- 002_create_refresh_tokens.up.sql
-- Sprint 1, GRAM-15: refresh tokens table for rotating JWT refresh strategy.
-- Run as grama_admin. Requires migration 001 (users table) to have been applied.

CREATE TABLE refresh_tokens (
    id           UUID        PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id      UUID        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash   TEXT        NOT NULL,
    expires_at   TIMESTAMPTZ NOT NULL,
    revoked_at   TIMESTAMPTZ,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Unique hash index: one lookup per validation request.
CREATE UNIQUE INDEX idx_refresh_tokens_hash       ON refresh_tokens(token_hash);
-- User index: revoke all tokens for a user on password change or logout-all.
CREATE INDEX        idx_refresh_tokens_user_id    ON refresh_tokens(user_id);
-- Expiry index: cleanup job scans expired tokens efficiently.
CREATE INDEX        idx_refresh_tokens_expires_at ON refresh_tokens(expires_at);

-- Application role: can insert new tokens, read and mark as revoked. No DELETE, no DDL.
GRANT SELECT, INSERT, UPDATE ON refresh_tokens TO grama_user;
