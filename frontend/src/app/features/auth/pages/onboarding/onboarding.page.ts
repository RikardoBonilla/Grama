import { Component } from '@angular/core';
import { Router } from '@angular/router';
import { IonContent, IonButton } from '@ionic/angular/standalone';
import { NgIf, NgFor } from '@angular/common';
import { Preferences } from '@capacitor/preferences';

interface Slide {
  icon: string;
  title: string;
  description: string;
}

@Component({
  selector: 'app-onboarding',
  standalone: true,
  imports: [IonContent, IonButton, NgIf, NgFor],
  template: `
    <ion-content>
      <div class="onboarding-wrapper">
        <div class="slide" *ngFor="let slide of slides; let i = index" [hidden]="i !== current">
          <div class="slide-icon">{{ slide.icon }}</div>
          <h2>{{ slide.title }}</h2>
          <p>{{ slide.description }}</p>
        </div>

        <div class="dots">
          <span *ngFor="let slide of slides; let i = index"
                class="dot" [class.active]="i === current"></span>
        </div>

        <div class="actions">
          <ion-button *ngIf="current < slides.length - 1" (click)="next()" expand="block">
            Siguiente
          </ion-button>
          <ion-button *ngIf="current === slides.length - 1" (click)="start()" expand="block">
            Empezar
          </ion-button>
        </div>
      </div>
    </ion-content>
  `,
  styles: [`
    .onboarding-wrapper {
      display: flex; flex-direction: column; align-items: center;
      justify-content: center; height: 100%; padding: 2rem; text-align: center;
    }
    .slide-icon { font-size: 5rem; margin-bottom: 1rem; }
    h2 { font-size: 1.6rem; font-weight: 700; margin-bottom: 0.5rem; }
    p { color: var(--ion-color-medium); margin-bottom: 2rem; }
    .dots { display: flex; gap: 0.5rem; margin-bottom: 2rem; }
    .dot { width: 10px; height: 10px; border-radius: 50%; background: var(--ion-color-medium); }
    .dot.active { background: var(--ion-color-primary); }
    .actions { width: 100%; }
  `],
})
export class OnboardingPage {
  current = 0;

  readonly slides: Slide[] = [
    { icon: '⚽', title: 'Encuentra canchas cerca de ti', description: 'Descubre complejos deportivos en tu ciudad en segundos.' },
    { icon: '📅', title: 'Reserva en segundos', description: 'Elige tu horario y reserva sin filas, directo desde tu celular.' },
    { icon: '💳', title: 'Paga fácil y seguro', description: 'Múltiples medios de pago con seguridad garantizada.' },
  ];

  constructor(private router: Router) {}

  next(): void {
    if (this.current < this.slides.length - 1) this.current++;
  }

  async start(): Promise<void> {
    await Preferences.set({ key: 'onboarding_done', value: 'true' });
    this.router.navigate(['/login']);
  }
}
