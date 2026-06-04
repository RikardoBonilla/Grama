import { Component, OnInit } from '@angular/core';
import { Router } from '@angular/router';
import { IonContent } from '@ionic/angular/standalone';
import { Preferences } from '@capacitor/preferences';
import { AuthService } from '../../../../core/auth/auth.service';

@Component({
  selector: 'app-splash',
  standalone: true,
  imports: [IonContent],
  template: `
    <ion-content class="ion-text-center">
      <div class="splash-container">
        <h1 class="splash-logo">Grama</h1>
        <p>Tu cancha, a un toque.</p>
      </div>
    </ion-content>
  `,
  styles: [`
    .splash-container {
      display: flex;
      flex-direction: column;
      align-items: center;
      justify-content: center;
      height: 100%;
      animation: fadeIn 0.8s ease-in;
    }
    .splash-logo { font-size: 3rem; font-weight: 700; color: var(--ion-color-primary); }
    @keyframes fadeIn { from { opacity: 0; } to { opacity: 1; } }
  `],
})
export class SplashPage implements OnInit {
  constructor(private auth: AuthService, private router: Router) {}

  async ngOnInit(): Promise<void> {
    await new Promise((r) => setTimeout(r, 2000));
    const user = await this.auth.getCurrentUser();

    if (user) {
      const destinations: Record<string, string> = {
        owner: '/admin/dashboard',
        operator: '/operator/agenda',
        client: '/home',
      };
      this.router.navigate([destinations[user.role] ?? '/home']);
      return;
    }

    const { value: onboardingDone } = await Preferences.get({ key: 'onboarding_done' });
    this.router.navigate([onboardingDone ? '/login' : '/onboarding']);
  }
}
