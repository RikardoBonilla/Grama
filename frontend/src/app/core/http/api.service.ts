// TODO Sprint 1: ApiService — wrapper sobre HttpClient de Angular.
//
// Responsabilidades:
//   - Centraliza la URL base del backend (desde environment.apiBaseUrl, NUNCA hardcodeada).
//   - Tipado fuerte de respuestas — sin `any`.
//   - Manejo centralizado de errores HTTP (delega al ErrorInterceptor para los genéricos,
//     y propaga los errores de dominio 4xx al llamador para que los muestre al usuario).
