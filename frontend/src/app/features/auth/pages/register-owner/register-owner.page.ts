import { Component, inject } from '@angular/core';
import { Router } from '@angular/router';
import { AbstractControl, FormBuilder, ReactiveFormsModule, ValidationErrors, ValidatorFn, Validators } from '@angular/forms';
import {
  IonContent, IonHeader, IonToolbar, IonTitle,
  IonItem, IonInput, IonButton, IonText, IonCheckbox, IonLabel,
  IonSelect, IonSelectOption,
} from '@ionic/angular/standalone';
import { NgFor, NgIf } from '@angular/common';
import { AuthService } from '../../../../core/auth/auth.service';

const passwordsMatch: ValidatorFn = (g: AbstractControl): ValidationErrors | null => {
  const pw = g.get('password')?.value;
  const c = g.get('confirmPassword')?.value;
  return pw === c ? null : { passwordsMismatch: true };
};

const CITIES = ['Bogota','Medellin','Cali','Barranquilla','Bucaramanga',
                 'Pereira','Manizales','Cartagena','Cucuta','Ibague'];

@Component({
  selector: 'app-register-owner',
  standalone: true,
  imports: [
    ReactiveFormsModule, NgIf,
    IonContent, IonHeader, IonToolbar, IonTitle,
    IonItem, IonInput, IonButton, IonText, IonCheckbox, IonLabel,
    IonSelect, IonSelectOption, NgFor,
  ],
  template: `
    <ion-header>
      <ion-toolbar>
        <ion-title>Registrar complejo</ion-title>
      </ion-toolbar>
    </ion-header>
    <ion-content class="ion-padding">
      <form [formGroup]="form" (ngSubmit)="onSubmit()">
        <ion-item>
          <ion-input formControlName="name" label="Tu nombre" labelPlacement="floating"></ion-input>
        </ion-item>
        <ion-item>
          <ion-input formControlName="companyName" label="Nombre del complejo" labelPlacement="floating"></ion-input>
        </ion-item>
        <ion-item>
          <ion-input formControlName="nit" label="NIT (ej: 123456789-0)" labelPlacement="floating"></ion-input>
        </ion-item>
        <ion-item>
          <ion-select formControlName="city" label="Ciudad" labelPlacement="floating">
            <ion-select-option *ngFor="let c of cities" [value]="c">{{ c }}</ion-select-option>
          </ion-select>
        </ion-item>
        <ion-item>
          <ion-input formControlName="phone" type="tel" label="Telefono (10 digitos)" labelPlacement="floating"></ion-input>
        </ion-item>
        <ion-item>
          <ion-input formControlName="email" type="email" label="Correo electronico" labelPlacement="floating"></ion-input>
        </ion-item>
        <ion-item>
          <ion-input formControlName="password" type="password" label="Contrasena" labelPlacement="floating"></ion-input>
        </ion-item>
        <ion-item>
          <ion-input formControlName="confirmPassword" type="password" label="Confirmar contrasena" labelPlacement="floating"></ion-input>
        </ion-item>
        <ion-item lines="none" class="ion-margin-top">
          <ion-checkbox formControlName="consentGiven" slot="start"></ion-checkbox>
          <ion-label class="ion-text-wrap ion-padding-start">
            Acepto el tratamiento de datos personales segun la Ley 1581 de 2012
          </ion-label>
        </ion-item>

        <ion-text color="danger" *ngIf="errorMsg">
          <p class="ion-padding-start">{{ errorMsg }}</p>
        </ion-text>

        <ion-button type="submit" expand="block" class="ion-margin-top"
                    [disabled]="form.invalid || isLoading">
          Registrar
        </ion-button>
      </form>
    </ion-content>
  `,
})
export class RegisterOwnerPage {
  private readonly auth = inject(AuthService);
  private readonly router = inject(Router);
  private readonly fb = inject(FormBuilder);

  isLoading = false;
  errorMsg = '';
  readonly cities = CITIES;

  readonly form = this.fb.nonNullable.group(
    {
      name: ['', Validators.required],
      companyName: ['', Validators.required],
      nit: ['', [Validators.required, Validators.pattern(/^\d{9}-\d$/)]],
      city: ['', Validators.required],
      phone: ['', [Validators.required, Validators.pattern(/^[0-9]{10}$/)]],
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
    // companyName, nit, city are persisted in Sprint 2 when creating the Venue.
    this.auth.register({ name, email, password, role: 'owner', consent_given: true }).subscribe({
      next: () => {
        this.isLoading = false;
        this.router.navigate(['/admin/dashboard']);
      },
      error: () => {
        this.isLoading = false;
        this.errorMsg = 'No pudimos crear tu cuenta. Intenta de nuevo.';
      },
    });
  }
}
