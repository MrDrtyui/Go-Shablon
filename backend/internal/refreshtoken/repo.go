package refreshtoken

import (
	"app/ent"
	"app/ent/refreshtoken"
	"app/internal/db"
	"context"
	"errors"
	"time"
)

var ErrRefreshTokenNotFound = errors.New("refresh token not found")

// Repository defines the interface for refresh token data access.
type Repository interface {
	Create(ctx context.Context, userID int, tokenHash string, expiresAt time.Time) (*ent.RefreshToken, error)
	GetByTokenHash(ctx context.Context, tokenHash string) (*ent.RefreshToken, error)
	Revoke(ctx context.Context, tokenHash string) error
	RevokeAllForUser(ctx context.Context, userID int) error
}

// PostgresRepo implements Repository using PostgreSQL via Ent.
type PostgresRepo struct {
	Db *db.Db
}

// NewPostgresRepo creates a new PostgreSQL repository.
func NewPostgresRepo(db *db.Db) Repository {
	return &PostgresRepo{Db: db}
}

// Create inserts a new refresh token into the database.
func (r *PostgresRepo) Create(ctx context.Context, userID int, tokenHash string, expiresAt time.Time) (*ent.RefreshToken, error) {
	return r.Db.Client.RefreshToken.Create().
		SetUserID(userID).
		SetTokenHash(tokenHash).
		SetExpiresAt(expiresAt).
		Save(ctx)
}

// GetByTokenHash retrieves a refresh token by its hash.
func (r *PostgresRepo) GetByTokenHash(ctx context.Context, tokenHash string) (*ent.RefreshToken, error) {
	token, err := r.Db.Client.RefreshToken.Query().
		Where(refreshtoken.TokenHashEQ(tokenHash)).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, ErrRefreshTokenNotFound
		}
		return nil, err
	}
	return token, nil
}

// Revoke marks a refresh token as revoked.
func (r *PostgresRepo) Revoke(ctx context.Context, tokenHash string) error {
	return r.Db.Client.RefreshToken.Update().
		Where(refreshtoken.TokenHashEQ(tokenHash)).
		SetRevoked(true).
		Exec(ctx)
}

// RevokeAllForUser marks all refresh tokens for a user as revoked.
func (r *PostgresRepo) RevokeAllForUser(ctx context.Context, userID int) error {
	return r.Db.Client.RefreshToken.Update().
		Where(refreshtoken.UserIDEQ(userID)).
		SetRevoked(true).
		Exec(ctx)
}
