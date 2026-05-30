// TODO Sprint 1: RoleGuard — protege rutas según el rol del usuario autenticado.
//
// Comportamiento:
//   - Lee el claim `role` del JWT del AuthService (nunca del JWT crudo del storage).
//   - Si el rol coincide con los permitidos en la ruta → permite la navegación.
//   - Si no → redirige a /unauthorized (403 page).
//
// Uso en el router:
//   { path: 'admin', component: AdminPage, canActivate: [authGuard, roleGuard],
//     data: { roles: ['owner', 'operator'] } }
//
// Seguridad: este guard es solo UX — la AUTORIZACIÓN REAL siempre ocurre en el backend.
//   Un atacante puede bypassear guards del frontend; el backend debe verificar el rol en cada request.
