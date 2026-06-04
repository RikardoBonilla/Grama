package postgres

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	domainuser "github.com/praxisvr/grama/internal/domain/user"
)

type PostgresUserRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresUserRepository(pool *pgxpool.Pool) *PostgresUserRepository {
	return &PostgresUserRepository{pool: pool}
}

const insertUser = `
INSERT INTO users
  (id, name, email, password_hash, role,
   consent_given, consent_given_at, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

func (r *PostgresUserRepository) Create(ctx context.Context, u *domainuser.User) error {
	// Store the raw hash bytes as hex so the TEXT column never holds binary garbage.
	hashHex := hex.EncodeToString(u.PasswordHash())

	_, err := r.pool.Exec(ctx, insertUser,
		u.ID().String(),
		u.Name(),
		u.Email(),
		hashHex,
		u.Role(),
		u.ConsentGiven(),
		u.ConsentGivenAt(),
		u.CreatedAt(),
		u.UpdatedAt(),
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			// Unique violation on email — return generic error to avoid
			// leaking whether a given email is already registered.
			return fmt.Errorf("registration failed")
		}
		return fmt.Errorf("create user: %w", err)
	}
	return nil
}

const selectUserCols = `
SELECT id, name, email, password_hash, role,
       consent_given, consent_given_at, created_at, updated_at
FROM users WHERE `

func (r *PostgresUserRepository) FindByEmail(ctx context.Context, email string) (*domainuser.User, error) {
	return r.scanUser(ctx, selectUserCols+"email = $1", email)
}

func (r *PostgresUserRepository) FindByID(ctx context.Context, id string) (*domainuser.User, error) {
	return r.scanUser(ctx, selectUserCols+"id = $1", id)
}

func (r *PostgresUserRepository) UpdateName(ctx context.Context, id, name string) error {
	_, err := r.pool.Exec(ctx,
		"UPDATE users SET name = $1 WHERE id = $2",
		name, id,
	)
	if err != nil {
		return fmt.Errorf("update user name: %w", err)
	}
	return nil
}

func (r *PostgresUserRepository) FindByRole(ctx context.Context, role string) ([]*domainuser.User, error) {
	rows, err := r.pool.Query(ctx, selectUserCols+"role = $1", role)
	if err != nil {
		return nil, fmt.Errorf("find by role: %w", err)
	}
	defer rows.Close()

	var users []*domainuser.User
	for rows.Next() {
		var (
			id, name, email, passwordHashHex, rowRole string
			consentGiven                               bool
			consentGivenAt                             *time.Time
			createdAt, updatedAt                       time.Time
		)
		if err := rows.Scan(
			&id, &name, &email, &passwordHashHex, &rowRole,
			&consentGiven, &consentGivenAt, &createdAt, &updatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan user row: %w", err)
		}
		var consentAt time.Time
		if consentGivenAt != nil {
			consentAt = *consentGivenAt
		}
		users = append(users, domainuser.Reconstitute(
			id, name, email, passwordHashHex, rowRole,
			consentGiven, consentAt, createdAt, updatedAt,
		))
	}
	return users, rows.Err()
}

func (r *PostgresUserRepository) scanUser(ctx context.Context, query string, arg any) (*domainuser.User, error) {
	row := r.pool.QueryRow(ctx, query, arg)

	var (
		id, name, email, passwordHashHex, role string
		consentGiven                            bool
		consentGivenAt                          *time.Time
		createdAt, updatedAt                    time.Time
	)

	err := row.Scan(
		&id, &name, &email, &passwordHashHex, &role,
		&consentGiven, &consentGivenAt, &createdAt, &updatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domainuser.ErrUserNotFound
		}
		return nil, fmt.Errorf("scan user: %w", err)
	}

	var consentAt time.Time
	if consentGivenAt != nil {
		consentAt = *consentGivenAt
	}

	return domainuser.Reconstitute(
		id, name, email, passwordHashHex, role,
		consentGiven,
		consentAt, createdAt, updatedAt,
	), nil
}
