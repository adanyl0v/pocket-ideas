package postgres

import (
	"context"
	"errors"
	"github.com/adanyl0v/pocket-ideas/internal/domain"
	"github.com/adanyl0v/pocket-ideas/internal/repository"
	"github.com/adanyl0v/pocket-ideas/pkg/database"
	"github.com/adanyl0v/pocket-ideas/pkg/log"
	"github.com/adanyl0v/pocket-ideas/pkg/proxerr"
	"github.com/adanyl0v/pocket-ideas/pkg/uuid"
	"time"
)

var (
	ErrUserAlreadyExists       = errors.New("user already exists")
	ErrUserFieldMustNotBeEmpty = errors.New("user field must not be empty")
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
	if tx == nil || conn == nil || !ok {
		r.logger.With(log.Fields{"tx": tx}).Error("failed to cast the transaction")
		return nil
	}

	return NewUserRepository(conn, r.logger, r.idGen)
}

const qInsertUser = `
INSERT INTO users (id, name, email, password, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6)
`

func (r *UserRepository) Save(ctx context.Context, user *domain.User) error {
	dto := newSaveUserDto(user)

	var err error
	dto.ID, err = r.idGen.NewV7()
	if err != nil {
		r.logger.WithError(err).Error("failed to generate user uuid")
		return err
	}

	dto.CreatedAt = time.Now()
	dto.UpdatedAt = dto.CreatedAt
	_, err = r.conn.Execute(ctx, qInsertUser, dto.ID, user.Name, dto.Email, dto.Password, dto.CreatedAt, dto.UpdatedAt)
	if err != nil {
		var pxErr proxerr.Error
		if errors.As(err, &pxErr) {
			e := pxErr.Unwrap()
			switch {
			case errors.Is(e, database.ErrUniqueViolation):
				err = proxerr.New(ErrUserAlreadyExists, pxErr.Error())
			case errors.Is(e, database.ErrNotNullViolation):
				err = proxerr.New(ErrUserFieldMustNotBeEmpty, pxErr.Error())
			}
		}

		r.logger.WithError(err).Error("failed to save a user")
		return err
	}

	r.logger.With(log.Fields{"id": dto.ID}).Debug("saved a user")
	dto.ToDomain(user)
	return nil
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
