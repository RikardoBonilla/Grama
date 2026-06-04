import { Component } from '@angular/core';
import { IonContent, IonHeader, IonToolbar, IonTitle } from '@ionic/angular/standalone';

@Component({
  selector: 'app-forbidden',
  standalone: true,
  imports: [IonContent, IonHeader, IonToolbar, IonTitle],
  template: `
    <ion-header><ion-toolbar><ion-title>Acceso denegado</ion-title></ion-toolbar></ion-header>
    <ion-content class="ion-padding ion-text-center">
      <h1>403</h1>
      <p>No tienes permiso para acceder a esta sección.</p>
    </ion-content>
  `,
})
export class ForbiddenPage {}
