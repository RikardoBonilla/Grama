# Multi-stage build — por qué esta estrategia:
# - Stage 1 (builder): imagen completa con el toolchain de Go para compilar.
# - Stage 2 (final): imagen "scratch" (vacía), solo el binario compilado.
# Resultado: imagen sin shell, sin OS tools, sin superficie de ataque.
# El proceso corre como "nobody" (UID 65534) — nunca como root.

# Stage 1: compilar
FROM golang:1.26-alpine AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# CGO_ENABLED=0 → binario estático, sin dependencias de C
# -trimpath → elimina rutas locales del binario (no filtra info del desarrollador)
# -ldflags="-s -w" → elimina símbolos de debug (reduce tamaño y superficie de info)
RUN CGO_ENABLED=0 GOOS=linux go build \
    -trimpath \
    -ldflags="-s -w" \
    -o /build/grama-api \
    ./cmd/api

# Stage 2: imagen final mínima
FROM scratch

# Certificados TLS para llamadas HTTPS salientes (ej: pasarela de pagos futura)
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copiar el binario
COPY --from=builder /build/grama-api /grama-api

# Correr como nobody (UID 65534) — nunca como root
USER 65534:65534

EXPOSE 8080

ENTRYPOINT ["/grama-api"]
