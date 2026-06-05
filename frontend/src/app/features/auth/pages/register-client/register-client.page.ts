import { Component, inject } from '@angular/core';
import { Router, RouterLink } from '@angular/router';
import { AbstractControl, FormBuilder, ReactiveFormsModule, ValidationErrors, ValidatorFn, Validators } from '@angular/forms';
import {
  IonContent, IonHeader, IonToolbar, IonTitle,
  IonItem, IonInput, IonButton, IonText, IonCheckbox, IonLabel,
} from '@ionic/angular/standalone';
import { NgIf } from '@angular/common';
import { AuthService } from '../../../../core/auth/auth.service';

const passwordsMatch: ValidatorFn = (group: AbstractControl): ValidationErrors | null => {
  const pw = group.get('password')?.value;
  const confirm = group.get('confirmPassword')?.value;
  return pw === confirm ? null : { passwordsMismatch: true };
};

@Component({
  selector: 'app-register-client',
  standalone: true,
  imports: [
    ReactiveFormsModule, RouterLink, NgIf,
    IonContent, IonHeader, IonToolbar, IonTitle,
    IonItem, IonInput, IonButton, IonText, IonCheckbox, IonLabel,
  ],
  template: `
    <ion-header>
      <ion-toolbar>
        <ion-title>Crear cuenta</ion-title>
      </ion-toolbar>
    </ion-header>
    <ion-content class="ion-padding">
      <form [formGroup]="form" (ngSubmit)="onSubmit()">
        <ion-item>
          <ion-input formControlName="name" label="Nombre completo"
                     labelPlacement="floating" autocomplete="name"></ion-input>
        </ion-item>
        <ion-item>
          <ion-input formControlName="email" type="email"
                     label="Correo electronico" labelPlacement="floating"></ion-input>
        </ion-item>
        <ion-item>
          <ion-input formControlName="password" type="password"
                     label="Contrasena" labelPlacement="floating"></ion-input>
        </ion-item>
        <ion-item>
          <ion-input formControlName="confirmPassword" type="password"
                     label="Confirmar contrasena" labelPlacement="floating"></ion-input>
        </ion-item>

        <ion-item lines="none" class="ion-margin-top">
          <ion-checkbox formControlName="consentGiven" slot="start"></ion-checkbox>
          <ion-label class="ion-text-wrap ion-padding-start">
            Acepto el tratamiento de mis datos personales segun la Ley 1581 de 2012
          </ion-label>
        </ion-item>

        <ion-text color="danger" *ngIf="errorMsg">
          <p class="ion-padding-start">{{ errorMsg }}</p>
        </ion-text>

        <ion-button type="submit" expand="block" class="ion-margin-top"
                    [disabled]="form.invalid || isLoading">
          Crear cuenta
        </ion-button>
      </form>

      <p class="ion-text-center ion-margin-top">
        Ya tienes cuenta? <a routerLink="/login">Ingresar</a>
      </p>
    </ion-content>
  `,
})
export class RegisterClientPage {
  private readonly auth = inject(AuthService);
  private readonly router = inject(Router);
  private readonly fb = inject(FormBuilder);

  isLoading = false;
  errorMsg = '';

  readonly form = this.fb.nonNullable.group(
    {
      name: ['', [Validators.required, Validators.minLength(2)]],
      email: ['', [Validators.required, Validators.email]],
      password: ['', [Validators.required, Validators.minLength(8)]],
      confirmPassword: ['', Validators.required],
      consentGiven: [false, Validators.requiredTrue],
    },
    { validators: passwordsMatch }
  );

  onSubmit(): void {
    if (this.form.invalid) return;
    this.isLoading = true;
    this.errorMsg = '';

    const { name, email, password } = this.form.getRawValue();
    this.auth.register({ name, email, password, role: 'client', consent_given: true }).subscribe({
      next: () => {
        this.isLoading = false;
        this.router.navigate(['/home']);
      },
      error: () => {
        this.isLoading = false;
        this.errorMsg = 'No pudimos crear tu cuenta. Intenta de nuevo.';
      },
    });
  }
}
