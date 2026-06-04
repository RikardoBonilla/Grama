package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/praxisvr/grama/internal/ports/repository"
)

type PostgresTokenRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresTokenRepository(pool *pgxpool.Pool) *PostgresTokenRepository {
	return &PostgresTokenRepository{pool: pool}
}

func (r *PostgresTokenRepository) SaveRefreshToken(
	ctx context.Context, userID, tokenHash string, expiresAt time.Time,
) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO refresh_tokens (user_id, token_hash, expires_at)
		VALUES ($1, $2, $3)`,
		userID, tokenHash, expiresAt,
	)
	if err != nil {
		return fmt.Errorf("save refresh token: %w", err)
	}
	return nil
}

func (r *PostgresTokenRepository) FindRefreshToken(
	ctx context.Context, tokenHash string,
) (*repository.RefreshToken, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT id, user_id, token_hash, expires_at, revoked_at, created_at
		FROM refresh_tokens
		WHERE token_hash = $1 AND revoked_at IS NULL`,
		tokenHash,
	)

	var t repository.RefreshToken
	var revokedAt *time.Time
	err := row.Scan(
		&t.ID, &t.UserID, &t.TokenHash,
		&t.ExpiresAt, &revokedAt, &t.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("refresh token not found")
		}
		return nil, fmt.Errorf("find refresh token: %w", err)
	}
	t.RevokedAt = revokedAt
	return &t, nil
}

func (r *PostgresTokenRepository) RevokeRefreshToken(
	ctx context.Context, tokenHash string,
) error {
	_, err := r.pool.Exec(ctx, `
		UPDATE refresh_tokens SET revoked_at = NOW()
		WHERE token_hash = $1`,
		tokenHash,
	)
	if err != nil {
		return fmt.Errorf("revoke refresh token: %w", err)
	}
	return nil
}

func (r *PostgresTokenRepository) RevokeAllUserTokens(
	ctx context.Context, userID string,
) error {
	_, err := r.pool.Exec(ctx, `
		UPDATE refresh_tokens SET revoked_at = NOW()
		WHERE user_id = $1`,
		userID,
	)
	if err != nil {
		return fmt.Errorf("revoke all user tokens: %w", err)
	}
	return nil
}
