-- init.sql — Inicialización de la base de datos Grama
-- Este archivo se ejecuta UNA sola vez cuando el contenedor se crea por primera vez.
-- Las migraciones incrementales van en migrations/.

-- TODO Sprint 0: implementar este archivo completo.
-- Contenido esperado:
--   1. Extensiones requeridas.
--   2. Rol grama_user con permisos mínimos.
--   3. Trigger update_updated_at_column() reutilizable.
--   4. Tabla audit_log para operaciones sensibles.

-- TODO: extensiones (ejecutar como superuser antes de crear grama_user)
-- CREATE EXTENSION IF NOT EXISTS "uuid-ossp";   -- uuid_generate_v4()
-- CREATE EXTENSION IF NOT EXISTS "pgcrypto";    -- gen_random_bytes(), crypt()
-- CREATE EXTENSION IF NOT EXISTS "btree_gist";  -- EXCLUDE constraint en reservas

-- TODO: crear rol de la aplicación con permisos mínimos
-- CREATE ROLE grama_user WITH LOGIN PASSWORD 'CHANGE_ME';
-- GRANT CONNECT ON DATABASE grama_db TO grama_user;
-- GRANT USAGE ON SCHEMA public TO grama_user;
-- GRANT SELECT, INSERT, UPDATE ON ALL TABLES IN SCHEMA public TO grama_user;
-- NOTA: sin DELETE, sin DDL, sin TRUNCATE — principio de mínimo privilegio.

-- TODO: función trigger para updated_at automático
-- CREATE OR REPLACE FUNCTION update_updated_at_column() ...

-- TODO: tabla audit_log
-- CREATE TABLE audit_log (
--   id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
--   table_name  TEXT NOT NULL,
--   record_id   UUID NOT NULL,
--   action      TEXT NOT NULL CHECK (action IN ('INSERT', 'UPDATE', 'DELETE')),
--   actor_id    UUID,
--   old_data    JSONB,
--   new_data    JSONB,
--   created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
-- );
