// TODO Sprint 1: JwtInterceptor — HttpInterceptor de Angular.
//
// Responsabilidades:
//   - Añade el header Authorization: Bearer <access_token> a todas las requests
//     que van a la API de Grama (no a dominios externos — verificar la URL).
//   - Si recibe 401, intenta refrescar el token UNA VEZ y reintenta la request original.
//   - Si el refresh también falla, redirige al login.
//
// Por qué en un interceptor y no en cada service:
//   DRY — un solo punto de gestión del token en toda la app.
//   Si cambiamos el esquema de auth, solo cambia aquí.
