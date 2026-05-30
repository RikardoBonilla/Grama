package http

// TODO Sprint 1: definir el router y montar todos los middlewares globales.
// Orden de middlewares (importa):
//   1. security.go (headers siempre primero — aplican a todas las respuestas, incluidos errores)
//   2. ratelimit.go (rechazar antes de parsear el token)
//   3. auth.go (validar JWT solo en rutas protegidas)
//   4. handlers
