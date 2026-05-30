# Multi-stage build para el frontend Ionic + Angular
#
# Por qué dos stages:
# - Stage 1 (builder): Node.js para compilar el proyecto Angular y producir el /www estático.
# - Stage 2 (server):  Nginx mínimo solo para servir los archivos estáticos compilados.
#   Resultado: imagen sin Node.js ni herramientas de build en producción → superficie mínima.

# Stage 1: compilar la app Angular
FROM node:22-alpine AS builder

WORKDIR /app

# Instalar dependencias primero (aprovecha el cache de Docker si package.json no cambia)
COPY package*.json ./
RUN npm ci --omit=dev=false

COPY . .

# Compilar para producción (output en /app/www)
RUN npm run build

# Stage 2: servir con Nginx
FROM nginx:1.27-alpine

# Copiar los archivos compilados
COPY --from=builder /app/www /usr/share/nginx/html

# Copiar configuración de Nginx para SPA (redirigir todo a index.html)
# TODO Sprint 1: crear nginx-frontend.conf con:
#   - Security headers (CSP, X-Frame-Options, etc.)
#   - try_files $uri /index.html (necesario para Angular routing)
#   - Caché agresivo para assets hasheados, no-cache para index.html
COPY nginx-frontend.conf /etc/nginx/conf.d/default.conf

EXPOSE 80
