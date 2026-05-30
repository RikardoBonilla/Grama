// TODO Sprint 1: AuthGuard — protege rutas que requieren estar autenticado.
//
// Comportamiento:
//   - Si el usuario tiene sesión activa → permite la navegación.
//   - Si no → redirige a /auth/login y guarda la URL intentada (returnUrl)
//     para redirigir después del login exitoso.
//
// Implementar como función (Angular 17+ usa functional guards, no clases CanActivate).
