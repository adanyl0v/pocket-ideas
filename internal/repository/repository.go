package repository

import (
	"context"
	"github.com/adanyl0v/pocket-ideas/internal/domain"
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetByID(ctx context.Context, id string) (domain.User, error)
	GetByEmail(ctx context.Context, email string) (domain.User, error)
	SelectAll(ctx context.Context) ([]domain.User, error)
	SelectByName(ctx context.Context, name string) ([]domain.User, error)
	UpdateByID(ctx context.Context, user *domain.User) error
	DeleteByID(ctx context.Context, id string) error
}
