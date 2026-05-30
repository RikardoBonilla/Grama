# Grama Frontend

Cliente web y móvil de Grama construido con Ionic + Angular + Capacitor.
Se compila para web, Android e iOS desde el mismo codebase.

## Prerrequisitos

- Node.js 22 LTS
- npm 10+
- `npm install -g @ionic/cli`
- Para Android: Android Studio + JDK 17
- Para iOS: Xcode 15+ (solo en macOS)

## Cómo correr en dev

```bash
cd frontend

# 1. Instalar dependencias (primera vez)
npm install

# 2. Copiar variables de entorno
cp .env.example .env.local
# Edita .env.local con la URL del backend local

# 3. Correr en el navegador
ionic serve

# 4. Correr en Android (requiere emulador o dispositivo conectado)
ionic capacitor run android -l

# 5. Correr en iOS (solo macOS)
ionic capacitor run ios -l
```

## Estructura

```
src/app/
├── core/          # Singletons: auth service, interceptores HTTP, guards
├── shared/        # Componentes y pipes reutilizables entre features
└── features/
    ├── auth/      # Login y registro (B2B y B2C)
    ├── venues/    # Listado público de complejos y canchas
    ├── bookings/  # Disponibilidad, reservar, cancelar
    └── admin/     # Panel B2B para dueños y operarios
```

## Comandos útiles

```bash
npm run type-check   # TypeScript strict — sin errores es requisito de mergear
npm run lint         # ESLint
npm run test         # Unit tests con coverage
npm run build        # Build de producción
```
