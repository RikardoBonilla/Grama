#!/usr/bin/env bash
# hardening.sh — Endurecer el VPS de producción (DigitalOcean Ubuntu 24 LTS)
# Uso: sudo ./scripts/hardening.sh
# TODO Sprint 4: implementar este script completo.
# Pasos esperados:
#   1. Actualizar el sistema (apt update && apt upgrade -y)
#   2. Configurar UFW: solo puertos 22, 80, 443
#   3. Deshabilitar login de root por SSH
#   4. Configurar fail2ban para SSH y Nginx
#   5. Instalar y configurar Certbot (TLS automático)
#   6. Configurar Docker para que Postgres NUNCA sea accesible externamente
#   7. Verificar que SonarQube solo escuche en localhost
#   8. Configurar rotación de logs
#   9. Habilitar actualizaciones de seguridad automáticas (unattended-upgrades)

set -euo pipefail
echo "TODO: implementar en Sprint 4"
