import { Component } from '@angular/core';
import { IonContent, IonHeader, IonToolbar, IonTitle } from '@ionic/angular/standalone';

@Component({
  selector: 'app-home',
  standalone: true,
  imports: [IonContent, IonHeader, IonToolbar, IonTitle],
  template: `
    <ion-header><ion-toolbar><ion-title>Home</ion-title></ion-toolbar></ion-header>
    <ion-content class="ion-padding"><p>Home -- coming soon (Sprint 2)</p></ion-content>
  `,
})
export class HomePage {}
