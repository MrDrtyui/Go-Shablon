package refreshtoken

import (
	"context"
	"errors"
	"time"
)

var (
	ErrInvalidRefreshToken = errors.New("invalid refresh token")
	ErrExpiredRefreshToken = errors.New("refresh token expired")
	ErrRevokedRefreshToken = errors.New("refresh token revoked")
)

// Service handles refresh token business logic.
type Service struct {
	Repo Repository
	TTL  time.Duration
}

// NewService creates a new refresh token service.
func NewService(repo Repository, ttl time.Duration) *Service {
	return &Service{
		Repo: repo,
		TTL:  ttl,
	}
}

// Generate creates a new refresh token for the given user.
// Returns the plain token string (to send to client) and any error.
func (s *Service) Generate(ctx context.Context, userID int) (string, error) {
	token, err := Generate()
	if err != nil {
		return "", err
	}

	tokenHash := Hash(token)
	expiresAt := time.Now().Add(s.TTL)

	_, err = s.Repo.Create(ctx, userID, tokenHash, expiresAt)
	if err != nil {
		return "", err
	}

	return token, nil
}

// Validate checks if the token is valid, not expired, and not revoked.
// Returns the associated userID if valid, or an error.
func (s *Service) Validate(ctx context.Context, token string) (int, error) {
	tokenHash := Hash(token)

	rt, err := s.Repo.GetByTokenHash(ctx, tokenHash)
	if err != nil {
		if errors.Is(err, ErrRefreshTokenNotFound) {
			return 0, ErrInvalidRefreshToken
		}
		return 0, err
	}

	if rt.Revoked {
		return 0, ErrRevokedRefreshToken
	}

	if time.Now().After(rt.ExpiresAt) {
		return 0, ErrExpiredRefreshToken
	}

	return rt.UserID, nil
}

// Rotate validates the old token, revokes it, and generates a new one.
// This implements single-use tokens with automatic rotation.
// Returns the new token, userID, and any error.
func (s *Service) Rotate(ctx context.Context, oldToken string) (string, int, error) {
	userID, err := s.Validate(ctx, oldToken)
	if err != nil {
		return "", 0, err
	}

	oldTokenHash := Hash(oldToken)
	if err := s.Repo.Revoke(ctx, oldTokenHash); err != nil {
		return "", 0, err
	}

	newToken, err := s.Generate(ctx, userID)
	if err != nil {
		return "", 0, err
	}

	return newToken, userID, nil
}

// Revoke marks the given token as revoked.
func (s *Service) Revoke(ctx context.Context, token string) error {
	tokenHash := Hash(token)
	return s.Repo.Revoke(ctx, tokenHash)
}
