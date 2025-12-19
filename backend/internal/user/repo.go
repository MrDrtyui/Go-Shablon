package user

import (
	"context"
	"errors"

	"app/ent"
	"app/ent/user"
	"app/internal/db"
)

var ErrUserNotFound = errors.New("User not found")

type Repository interface {
	Create(ctx context.Context, emailDto string, passwordHash string, username string) (*ent.User, error)
	GetByEmail(ctx context.Context, emailDto string) (*ent.User, error)
	GetById(ctx context.Context, id int) (*ent.User, error)
}

type PostgresRepo struct {
	Db *db.Db
}

func (p *PostgresRepo) Create(ctx context.Context, emailDto string, passwordHash string, username string) (*ent.User, error) {
	return p.Db.Client.User.Create().SetEmail(emailDto).SetPassword(passwordHash).SetNillableUsername(&username).Save(ctx)
}

func (p *PostgresRepo) GetByEmail(ctx context.Context, emailDto string) (*ent.User, error) {
	u, err := p.Db.Client.User.Query().Where(user.EmailEQ(emailDto)).Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, ErrUserNotFound
		}
		return nil, err

	}
	return u, nil
}

func (p *PostgresRepo) GetById(ctx context.Context, id int) (*ent.User, error) {
	u, err := p.Db.Client.User.Get(ctx, id)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return u, nil
}

func NewPostgresRepo(db *db.Db) Repository {
	return &PostgresRepo{Db: db}
}
