package postgres

import (
	"context"
	"github.com/adanyl0v/pocket-ideas/internal/domain"
	"github.com/adanyl0v/pocket-ideas/internal/repository"
	"github.com/adanyl0v/pocket-ideas/pkg/database"
	"github.com/adanyl0v/pocket-ideas/pkg/log"
	"github.com/adanyl0v/pocket-ideas/pkg/uuid"
)

type UserRepository struct {
	Repository
	idGen uuid.Generator
}

func NewUserRepository(conn database.Conn, logger log.Logger, idGen uuid.Generator) *UserRepository {
	return &UserRepository{
		Repository: Repository{
			conn:   conn,
			logger: logger,
		},
		idGen: idGen,
	}
}

func (r *UserRepository) WithTx(tx repository.Tx) repository.Repository {
	conn, ok := tx.(database.Tx)
	if !ok {
		r.logger.With(log.Fields{"tx": tx}).Error("failed to cast the transaction")
		return nil
	}

	return NewUserRepository(conn, r.logger, r.idGen)
}

func (r *UserRepository) Save(ctx context.Context, user *domain.User) error {
	// TODO implement me
	panic("implement me")
}

func (r *UserRepository) FindById(ctx context.Context, id string) (domain.User, error) {
	// TODO implement me
	panic("implement me")
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	// TODO implement me
	panic("implement me")
}

func (r *UserRepository) FindAll(ctx context.Context) ([]domain.User, error) {
	// TODO implement me
	panic("implement me")
}

func (r *UserRepository) FindByName(ctx context.Context, name string) ([]domain.User, error) {
	// TODO implement me
	panic("implement me")
}

func (r *UserRepository) UpdateById(ctx context.Context, user *domain.User) error {
	// TODO implement me
	panic("implement me")
}

func (r *UserRepository) DeleteById(ctx context.Context, id string) error {
	// TODO implement me
	panic("implement me")
}
