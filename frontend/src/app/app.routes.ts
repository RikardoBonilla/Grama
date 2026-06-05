import { Routes } from '@angular/router';
import { authGuard } from './core/guards/auth.guard';
import { roleGuard } from './core/guards/role.guard';

export const routes: Routes = [
  { path: '', redirectTo: 'splash', pathMatch: 'full' },
  {
    path: 'splash',
    loadComponent: () => import('./features/auth/pages/splash/splash.page').then(m => m.SplashPage),
  },
  {
    path: 'onboarding',
    loadComponent: () => import('./features/auth/pages/onboarding/onboarding.page').then(m => m.OnboardingPage),
  },
  {
    path: 'login',
    loadComponent: () => import('./features/auth/pages/login/login.page').then(m => m.LoginPage),
  },
  {
    path: 'register',
    children: [
      {
        path: 'client',
        loadComponent: () => import('./features/auth/pages/register-client/register-client.page').then(m => m.RegisterClientPage),
      },
      {
        path: 'owner',
        loadComponent: () => import('./features/auth/pages/register-owner/register-owner.page').then(m => m.RegisterOwnerPage),
      },
    ],
  },
  {
    path: 'home',
    canActivate: [authGuard],
    loadComponent: () => import('./features/home/home.page').then(m => m.HomePage),
  },
  {
    path: 'admin',
    canActivate: [authGuard, roleGuard],
    data: { role: 'owner' },
    loadChildren: () => import('./features/admin/admin.routes').then(m => m.ADMIN_ROUTES),
  },
  {
    path: 'operator',
    canActivate: [authGuard, roleGuard],
    data: { role: 'operator' },
    loadChildren: () => import('./features/operator/operator.routes').then(m => m.OPERATOR_ROUTES),
  },
  {
    path: '403',
    loadComponent: () => import('./features/shared/forbidden/forbidden.page').then(m => m.ForbiddenPage),
  },
  { path: '**', redirectTo: 'splash' },
];
