package repository

import (
	"context"
	"github.com/adanyl0v/pocket-ideas/internal/domain"
)

type UserRepository interface {
	Repository
	Create(ctx context.Context, user *domain.User) error
	GetByID(ctx context.Context, id string) (domain.User, error)
	GetByEmail(ctx context.Context, email string) (domain.User, error)
	SelectAll(ctx context.Context) ([]domain.User, error)
	SelectByName(ctx context.Context, name string) ([]domain.User, error)
	UpdateByID(ctx context.Context, user *domain.User) error
	DeleteByID(ctx context.Context, id string) error
}

// This is a simple workaround for adding a transaction support without the need
// to implement complex patterns, such as [UoW] and without breaking the existing
// architecture by encapsulating the database logic exclusively within a repository.
//
// [UoW]: https://en.wikipedia.org/wiki/Unit_of_work
type (
	Repository interface {
		// Begin a new implementation-specific transaction
		Begin(ctx context.Context) (Tx, error)

		// WithTx copies the repository and connects it to the transaction
		WithTx(tx Tx) Repository
	}

	Tx interface {
		Commit(ctx context.Context) error
		Rollback(ctx context.Context) error
	}
)
