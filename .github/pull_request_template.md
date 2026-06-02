## Descripción
<!-- Qué hace este PR y por qué -->

## Historia Jira
<!-- Clave GRAM-XX -->

## Checklist Definición de Hecho
- [ ] Código en feature branch, respeta reglas de arquitectura hexagonal
- [ ] Queries parametrizadas ($1, $2) — sin fmt.Sprintf con datos de usuario
- [ ] Nunca float64 para dinero — int64 en centavos
- [ ] Endpoint con auth + rate limiting si corresponde
- [ ] Autorización por rol Y por pertenencia verificada
- [ ] Tests unitarios escritos y en verde
- [ ] go vet ./... en verde
- [ ] go test -race ./... en verde
- [ ] gosec ./... sin HIGH/MEDIUM
- [ ] govulncheck ./... en verde
- [ ] SonarQube Quality Gate verde (sin Security Hotspots)
- [ ] Commit semántico con clave GRAM al final
- [ ] Historia marcada como Done en Jira

## Tipo de cambio
- [ ] feat — nueva funcionalidad
- [ ] fix — corrección de bug
- [ ] chore — infraestructura / configuración
- [ ] test — tests
- [ ] docs — documentación
