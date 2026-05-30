# Architecture Decision Records — Grama

## ADR-001: Clean Architecture + Hexagonal + DDD

**Estado**: Aceptado | **Sprint**: 0

**Contexto**: Necesitamos una arquitectura que sea testeable, segura y que permita cambiar
la base de datos o el framework HTTP sin reescribir la lógica de negocio.

**Decisión**: Monolito modular con Clean Architecture + Hexagonal Architecture + DDD.

**Consecuencias**:
- domain/ es el núcleo puro: entidades, value objects, errores de dominio. Sin dependencias externas.
- ports/ define las interfaces (contratos) que el dominio espera del mundo exterior.
- application/ orquesta casos de uso usando solo interfaces, nunca implementaciones.
- infrastructure/ contiene las implementaciones concretas (PostgreSQL, HTTP, JWT).

**Regla de dependencia**: las capas internas no conocen a las externas. La dependencia siempre
apunta hacia adentro.

---

## ADR-002: Anti-doble-reserva con defensa en profundidad

**Estado**: Aceptado | **Sprint**: 3

**Contexto**: Dos clientes pueden intentar reservar el mismo slot al mismo tiempo (race condition).

**Decisión**: Dos capas de defensa:
1. Validación en Go (application): chequeo optimista antes de insertar.
2. EXCLUDE constraint en PostgreSQL con btree_gist sobre (court_id, tstzrange) donde status != 'cancelled'.
   La base de datos rechaza físicamente el solapamiento aunque la validación de Go falle.

**Por qué no solo en Go**: una validación de aplicación puede perder ante concurrencia real
(dos requests pasan el check al mismo tiempo antes de que ninguno inserte).
El constraint de base de datos es atómico y transaccional.

---

## ADR-003: UUIDs v4 como identificadores

**Estado**: Aceptado | **Sprint**: 0

**Contexto**: Los IDs secuenciales (BIGSERIAL) filtran información de negocio (cuántas reservas
hay, cuántos usuarios) y permiten ataques de enumeración (IDOR trivial).

**Decisión**: UUID v4 generados con uuid_generate_v4() en PostgreSQL para todas las tablas.

**Consecuencias**: IDs no predecibles, no enumerables. La autorización por pertenencia sigue
siendo necesaria (un UUID no es un control de acceso).

---

## ADR-004: Montos en centavos (enteros), nunca float

**Estado**: Aceptado | **Sprint**: 4

**Contexto**: Los floats tienen errores de precisión que son inaceptables en operaciones de dinero.

**Decisión**: Todos los montos se almacenan y procesan como enteros (centavos de peso colombiano).
Ejemplo: $50.000 COP = 5000000 centavos.

**Consecuencias**: Evitamos errores de redondeo. La lógica de presentación (dividir entre 100)
ocurre solo en la capa de UI.
