-- 002_create_refresh_tokens.down.sql
-- Rolls back migration 002. CASCADE drops indexes and grants automatically.
-- Run as grama_admin.

DROP TABLE IF EXISTS refresh_tokens CASCADE;
