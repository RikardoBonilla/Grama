import { CanActivateFn, ActivatedRouteSnapshot, Router } from '@angular/router';
import { inject } from '@angular/core';
import { AuthService } from '../auth/auth.service';

export const roleGuard: CanActivateFn = async (route: ActivatedRouteSnapshot) => {
  const auth = inject(AuthService);
  const router = inject(Router);

  const user = await auth.getCurrentUser();
  if (!user) {
    router.navigate(['/login']);
    return false;
  }

  const required: string = route.data['role'] ?? '';
  if (required && user.role !== required) {
    router.navigate(['/403']);
    return false;
  }
  return true;
};
