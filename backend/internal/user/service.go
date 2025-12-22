package user

import (
	"app/ent"
	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrInvalidPassword    = errors.New("invalid password")
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

func (s *Service) Login(ctx context.Context, email string, password string) (*ent.User, error) {
	u, err := s.Repo.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	return u, nil
}

func (s *Service) GetByEmail(ctx context.Context, email string) (*ent.User, error) {
	return s.Repo.GetByEmail(ctx, email)
}

func (s *Service) GetByID(ctx context.Context, id int) (*ent.User, error) {
	return s.Repo.GetById(ctx, id)
}
