-- 001_create_users.down.sql
-- Rolls back migration 001. CASCADE drops indexes, triggers, and grants automatically.
-- Run as grama_admin.

DROP TABLE IF EXISTS users CASCADE;
