// TODO Sprint 4: variables de entorno para PRODUCCIÓN.
// La URL real del backend se inyecta en el build de CI vía variable de entorno.
// NUNCA hardcodear la URL de producción aquí — cambiaría entre deploys.

export const environment = {
  production: true,
  apiBaseUrl: 'https://api.grama.app/api/v1',  // TODO: reemplazar con dominio real en Sprint 4
};
