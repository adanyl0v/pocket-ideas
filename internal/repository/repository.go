package repository

import (
	"context"
	"github.com/adanyl0v/pocket-ideas/internal/domain"
)

// The simple workaround for adding a transaction support without the need to implement
// complex patterns, such as [UoW] and without breaking the "clean" architecture rules.
//
// [UoW]: https://en.wikipedia.org/wiki/Unit_of_work
type (
	Tx interface {
		Commit(ctx context.Context) error
		Rollback(ctx context.Context) error
	}

	Repository interface {
		// Begin an implementation-specific transaction
		Begin(ctx context.Context) (Tx, error)

		// WithTx copies the repository and sets up the underlying connection
		WithTx(tx Tx) Repository
	}
)

type UserRepository interface {
	Repository
	Save(ctx context.Context, user *domain.User) error
	FindById(ctx context.Context, id string) (domain.User, error)
	FindByEmail(ctx context.Context, email string) (domain.User, error)
	FindAll(ctx context.Context) ([]domain.User, error)
	FindByName(ctx context.Context, name string) ([]domain.User, error)
	UpdateById(ctx context.Context, user *domain.User) error
	DeleteById(ctx context.Context, id string) error
}
