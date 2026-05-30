# Grama

Plataforma web y móvil para la administración de complejos de canchas sintéticas en Colombia.

## Visión del producto

Grama permite a dueños de complejos gestionar sus canchas, franjas horarias y cobros.
A los jugadores les permite ver disponibilidad, reservar y pagar en pocos toques.

**El corazón del producto**: reservas sin doble-booking, resuelto con defensa en profundidad
(validación en Go + EXCLUDE constraint en PostgreSQL con btree_gist).

## Roles

| Rol | Mundo | Función |
|---|---|---|
| Dueño/Administrador | B2B | Gestiona el complejo, canchas, precios y reportes |
| Operario/Recepción | B2B | Crea reservas presenciales, registra pagos |
| Cliente/Jugador | B2C | Busca disponibilidad, reserva, cancela |

## Stack

- **Backend**: Go 1.22+, Clean Architecture + DDD + Hexagonal
- **Base de datos**: PostgreSQL 16 (pgx/v5, btree_gist, tstzrange)
- **Frontend**: Ionic + Angular (repositorio separado)
- **Infra**: Docker + Nginx + DigitalOcean VPS

## Cómo correr el proyecto en local

```bash
# 1. Copia las variables de entorno
cp .env.example .env
# Edita .env con tus valores locales

# 2. Levanta la base de datos
docker compose up -d grama_postgres

# 3. Corre el backend
cd backend
go run cmd/api/main.go
```

Consulta `docs/ARCHITECTURE.md` para decisiones de diseño y `docs/ROADMAP.md` para el estado del MVP.
