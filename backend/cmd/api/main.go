package main

// TODO: entrypoint principal del servidor HTTP.
// Responsabilidades de main.go (y SOLO de main.go):
//   1. Leer configuración del entorno.
//   2. Construir las dependencias concretas (DB pool, repositorios, casos de uso, handlers).
//   3. Conectar todo mediante inyección de dependencias.
//   4. Arrancar el servidor con graceful shutdown (señales SIGINT/SIGTERM).
//
// NUNCA debe haber lógica de negocio aquí. main.go es el "pegamento".
