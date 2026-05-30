#!/usr/bin/env bash
# setup.sh — Levantar entorno local de desarrollo
# Uso: ./scripts/setup.sh

set -euo pipefail

echo "==> [Grama] Verificando dependencias..."
command -v docker >/dev/null 2>&1 || { echo "ERROR: Docker no está instalado."; exit 1; }
command -v go >/dev/null 2>&1 || { echo "ERROR: Go no está instalado."; exit 1; }

echo "==> [Grama] Copiando .env.example a .env (si no existe)..."
if [ ! -f .env ]; then
  cp .env.example .env
  echo "IMPORTANTE: Edita .env con tus valores locales antes de continuar."
  exit 0
fi

echo "==> [Grama] Levantando PostgreSQL..."
docker compose up -d grama_postgres

echo "==> [Grama] Esperando a que PostgreSQL esté listo..."
sleep 3
docker compose exec grama_postgres pg_isready -U "${POSTGRES_USER:-postgres}" && echo "PostgreSQL OK"

echo "==> [Grama] Descargando dependencias de Go..."
cd backend && go mod download

echo ""
echo "✓ Entorno listo. Corre: cd backend && go run cmd/api/main.go"
