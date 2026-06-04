package auth

// TokenManager is the port for JWT operations.
// The application layer depends on this interface; infrastructure/auth implements it.
type TokenManager interface {
	GenerateAccessToken(userID, role string) (string, error)
}
