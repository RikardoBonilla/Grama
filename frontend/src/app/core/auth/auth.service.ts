import { Injectable, inject } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable, from, switchMap, tap } from 'rxjs';
import { Preferences } from '@capacitor/preferences';
import { Router } from '@angular/router';
import { environment } from '../../../environments/environment';

export interface LoginResponse {
  access_token: string;
  refresh_token: string;
  expires_in: number;
  user: { id: string; name: string; role: string };
}

export interface RegisterData {
  name: string;
  email: string;
  password: string;
  role: string;
  consent_given: boolean;
}

export interface RegisterResponse {
  id: string;
  name: string;
  email: string;
  role: string;
  created_at: string;
}

export interface TokenResponse {
  access_token: string;
  refresh_token: string;
  expires_in: number;
}

export interface UserInfo {
  id: string;
  name: string;
  role: 'owner' | 'operator' | 'client';
}

const REFRESH_KEY = 'refresh_token';

/** Decodes a JWT payload without verifying the signature (verification is backend's job). */
function decodeJwt(token: string): Record<string, unknown> | null {
  try {
    const payload = token.split('.')[1];
    return JSON.parse(atob(payload.replace(/-/g, '+').replace(/_/g, '/')));
  } catch {
    return null;
  }
}

@Injectable({ providedIn: 'root' })
export class AuthService {
  private readonly http = inject(HttpClient);
  private readonly router = inject(Router);
  private readonly base = environment.apiUrl;

  // Access token lives in memory only — XSS cannot reach it.
  private accessToken: string | null = null;

  login(email: string, password: string): Observable<LoginResponse> {
    return this.http.post<LoginResponse>(`${this.base}/auth/login`, { email, password }).pipe(
      tap(async (res) => {
        this.accessToken = res.access_token;
        await Preferences.set({ key: REFRESH_KEY, value: res.refresh_token });
      })
    );
  }

  register(data: RegisterData): Observable<RegisterResponse> {
    return this.http.post<RegisterResponse>(`${this.base}/auth/register`, data);
  }

  logout(): Observable<void> {
    return from(Preferences.get({ key: REFRESH_KEY })).pipe(
      switchMap(({ value }) =>
        this.http.post<void>(`${this.base}/auth/logout`, { refresh_token: value })
      ),
      tap(async () => {
        this.accessToken = null;
        await Preferences.remove({ key: REFRESH_KEY });
        this.router.navigate(['/login']);
      })
    );
  }

  refreshToken(): Observable<TokenResponse> {
    return from(Preferences.get({ key: REFRESH_KEY })).pipe(
      switchMap(({ value }) =>
        this.http.post<TokenResponse>(`${this.base}/auth/refresh`, { refresh_token: value })
      ),
      tap(async (res) => {
        this.accessToken = res.access_token;
        await Preferences.set({ key: REFRESH_KEY, value: res.refresh_token });
      })
    );
  }

  getAccessToken(): string | null {
    return this.accessToken;
  }

  async getCurrentUser(): Promise<UserInfo | null> {
    const { value: refreshToken } = await Preferences.get({ key: REFRESH_KEY });
    if (!refreshToken || !this.accessToken) return null;
    const payload = decodeJwt(this.accessToken);
    if (!payload) return null;
    return {
      id: payload['sub'] as string,
      name: (payload['name'] as string) ?? '',
      role: payload['role'] as 'owner' | 'operator' | 'client',
    };
  }

  async isAuthenticated(): Promise<boolean> {
    if (!this.accessToken) return false;
    const payload = decodeJwt(this.accessToken);
    if (!payload || typeof payload['exp'] !== 'number') return false;
    return payload['exp'] > Date.now() / 1000;
  }
}
