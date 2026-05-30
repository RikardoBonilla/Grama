// TODO Sprint 1: ErrorInterceptor — HttpInterceptor de Angular.
//
// Responsabilidades:
//   - 401 → delegar al JwtInterceptor (refresh + retry). Si falla → logout.
//   - 403 → mostrar mensaje "No tienes permiso para esta acción".
//   - 429 → mostrar mensaje "Demasiados intentos. Espera un momento." (rate limiting del backend).
//   - 5xx → mostrar mensaje genérico de error de servidor, loggear en consola (no en producción).
//
// Regla de seguridad: NUNCA mostrar al usuario el mensaje de error crudo del servidor.
//   Los mensajes de error de backend pueden revelar información de arquitectura interna.
