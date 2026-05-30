# Security Model — Grama

## Modelo de Amenazas (STRIDE simplificado)

| Amenaza | Vector | Control |
|---|---|---|
| Doble reserva | Concurrencia de requests | Validación en Go + EXCLUDE constraint PostgreSQL |
| IDOR | Acceso a recursos de otro usuario | Autorización por pertenencia en cada caso de uso |
| Inyección SQL | Input malicioso en queries | Queries parametrizadas SIEMPRE ($1, $2...) |
| Contraseña débil / texto plano | Brecha de DB | Argon2id con parámetros endurecidos |
| JWT inválido / alg:none | Tokens manipulados | Validación estricta de algoritmo en middleware |
| Enumeración de usuarios | Respuestas de error distintas | Mensajes de error genéricos siempre |
| Fuerza bruta | Login / registro / reset | Rate limiting por IP en endpoints sensibles |
| Información filtrada en headers | Fingerprinting del servidor | Security headers en todas las respuestas |
| Secrets en git | Commit accidental de .env | .gitignore + hook pre-commit anti-secretos |
| Vulnerabilidades en deps | Supply chain | govulncheck en CI + Trivy en imagen Docker |
| Datos personales sin protección | Ley 1581/2012 Colombia | Habeas Data: consentimiento, minimización, audit_log |

## Controles de autenticación

- **Contraseñas**: Argon2id (memory=64MB, iterations=3, parallelism=2, salt=16B aleatorio, key=32B)
- **JWT Access Token**: TTL 15 minutos, algoritmo HS256, validación estricta de alg
- **JWT Refresh Token**: TTL 7 días, rotativo (el viejo se invalida al renovar)
- **Bloqueo alg:none**: middleware rechaza tokens con algoritmo "none" — ataque conocido

## Controles de autorización

- **Por rol**: owner, operator, client — definido en el JWT claim
- **Por pertenencia**: cada caso de uso verifica que el recurso pertenece al usuario autenticado
  - Un operator solo puede operar sobre su complejo asignado
  - Un client solo puede ver/cancelar SUS reservas

## Red y contenedores

- PostgreSQL: bind en 127.0.0.1:5433, nunca 0.0.0.0
- Docker: imagen scratch, usuario nobody (UID 65534), sin shell, sin herramientas de debug
- Nginx: TLS 1.2/1.3 únicamente, HSTS, cabeceras de seguridad
- VPS: UFW solo puertos 22, 80, 443

## Cumplimiento Ley 1581 de 2012 (Habeas Data — Colombia)

- Se registran datos personales de clientes externos (no solo empleados).
- audit_log registra: quién creó/modificó/canceló qué reserva y cuándo.
- Consentimiento explícito en el registro de clientes.
- Minimización de datos: solo recopilamos lo necesario para la reserva.
- TODO Sprint 4: definir política de retención y proceso de eliminación de datos.
