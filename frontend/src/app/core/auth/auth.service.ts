// TODO Sprint 1: AuthService — singleton inyectado en root.
//
// Responsabilidades:
//   - login(email, password): llama al backend, recibe access_token + refresh_token.
//   - logout(): invalida el refresh_token en el backend y limpia el estado local.
//   - refreshToken(): rota el refresh_token silenciosamente antes de que expire el access_token.
//   - currentUser$: Observable<User | null> — el resto de la app consume este observable.
//
// Decisión de seguridad — dónde guardar los tokens:
//   - access_token: en MEMORIA (variable privada en el servicio). No sobrevive refresh de página,
//     pero eso está bien porque el refresh_token lo regenera automáticamente.
//   - refresh_token: en cookie HttpOnly + SameSite=Strict, servida por el backend.
//     El JavaScript de la app NUNCA puede leer una cookie HttpOnly → XSS no puede robarlo.
//   - NUNCA en localStorage: persiste entre sesiones y es accesible por cualquier script XSS.
//
// La decisión final sobre el mecanismo de storage se documenta en Sprint 1 tras evaluar
// las restricciones de Capacitor (apps móviles nativas tienen contexto de cookies diferente).
