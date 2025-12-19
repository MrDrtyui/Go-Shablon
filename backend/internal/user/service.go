package user

import (
	"app/ent"
	"context"

	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	Repo Repository
}

func NewSercie(repo Repository) *Service {
	return &Service{Repo: repo}
}

func (s *Service) Register(ctx context.Context, email string, password string, username string) (*ent.User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	u, err := s.Repo.Create(ctx, email, string(hash), username)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (s *Service) GetByEmail(ctx context.Context, email string) (*ent.User, error) {
	return s.Repo.GetByEmail(ctx, email)
}
