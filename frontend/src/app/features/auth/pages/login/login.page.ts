import { Component, inject } from '@angular/core';
import { Router, RouterLink } from '@angular/router';
import { FormBuilder, ReactiveFormsModule, Validators } from '@angular/forms';
import {
  IonContent, IonHeader, IonToolbar, IonTitle,
  IonItem, IonInput, IonButton, IonText, IonSpinner,
} from '@ionic/angular/standalone';
import { NgIf } from '@angular/common';
import { AuthService } from '../../../../core/auth/auth.service';

@Component({
  selector: 'app-login',
  standalone: true,
  imports: [
    ReactiveFormsModule, RouterLink, NgIf,
    IonContent, IonHeader, IonToolbar, IonTitle,
    IonItem, IonInput, IonButton, IonText, IonSpinner,
  ],
  template: `
    <ion-header>
      <ion-toolbar>
        <ion-title>Ingresar</ion-title>
      </ion-toolbar>
    </ion-header>
    <ion-content class="ion-padding">
      <form [formGroup]="form" (ngSubmit)="onSubmit()">
        <ion-item>
          <ion-input formControlName="email" type="email"
                     label="Correo electronico" labelPlacement="floating"
                     autocomplete="email"></ion-input>
        </ion-item>
        <ion-item>
          <ion-input formControlName="password" type="password"
                     label="Contrasena" labelPlacement="floating"
                     autocomplete="current-password"></ion-input>
        </ion-item>

        <ion-text color="danger" *ngIf="errorMsg">
          <p class="ion-padding-start">{{ errorMsg }}</p>
        </ion-text>

        <ion-button type="submit" expand="block" class="ion-margin-top"
                    [disabled]="form.invalid || isLoading">
          <ion-spinner *ngIf="isLoading" name="crescent"></ion-spinner>
          <span *ngIf="!isLoading">Ingresar</span>
        </ion-button>
      </form>

      <p class="ion-text-center ion-margin-top">
        No tienes cuenta?
        <a routerLink="/register/client">Registrate</a>
      </p>
    </ion-content>
  `,
})
export class LoginPage {
  private readonly auth = inject(AuthService);
  private readonly router = inject(Router);
  private readonly fb = inject(FormBuilder);

  isLoading = false;
  errorMsg = '';

  readonly form = this.fb.nonNullable.group({
    email: ['', [Validators.required, Validators.email]],
    password: ['', [Validators.required, Validators.minLength(8)]],
  });

  onSubmit(): void {
    if (this.form.invalid) return;
    this.isLoading = true;
    this.errorMsg = '';

    const { email, password } = this.form.getRawValue();
    this.auth.login(email, password).subscribe({
      next: (res) => {
        this.isLoading = false;
        const destinations: Record<string, string> = {
          owner: '/admin/dashboard',
          operator: '/operator/agenda',
          client: '/home',
        };
        this.router.navigate([destinations[res.user.role] ?? '/home']);
      },
      error: () => {
        this.isLoading = false;
        // Generic message — never reveal whether the email exists.
        this.errorMsg = 'Credenciales incorrectas. Intenta de nuevo.';
      },
    });
  }
}
