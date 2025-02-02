package postgres

import (
	"context"
	"errors"
	"github.com/adanyl0v/pocket-ideas/internal/domain"
	pgdb "github.com/adanyl0v/pocket-ideas/pkg/database/postgres"
	"github.com/adanyl0v/pocket-ideas/pkg/log"
	"github.com/adanyl0v/pocket-ideas/pkg/proxerr"
	uuidgen "github.com/adanyl0v/pocket-ideas/pkg/uuid"
	"time"
)

var (
	ErrUserExists = errors.New("user already exists")
)

type UserRepository struct {
	repository
	idGen uuidgen.Generator
}

func NewUserRepository(conn pgdb.Conn, logger log.Logger, idGen uuidgen.Generator) *UserRepository {
	return &UserRepository{
		repository: repository{
			conn:   conn,
			logger: logger,
		},
		idGen: idGen,
	}
}

const insertUserQuery = `
INSERT INTO users (id, name, email, password, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6)
`

func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	var dto createUserDTO
	dto.FromDomain(user)

	id, err := r.idGen.NewV7()
	if err != nil {
		r.logger.WithError(err).Error("failed to generate a new user id")
		return err
	}

	dto.ID = id
	dto.CreatedAt = time.Now()
	dto.UpdatedAt = dto.CreatedAt
	if err = r.conn.Exec(ctx, insertUserQuery, dto.ID, dto.Name, dto.Email,
		dto.Password, dto.CreatedAt, dto.UpdatedAt); err != nil {

		var pxErr proxerr.Error
		if errors.As(err, &pxErr) && errors.Is(pxErr.Unwrap(), pgdb.ErrUniqueViolation) {
			r.logger.With(log.Fields{
				"email": dto.Email,
			}).Warn("tried to create a new user with existing email")
			err = proxerr.New(ErrUserExists, pxErr.Error())
		} else {
			r.logger.WithError(err).Error("failed to insert user into the database")
		}

		return err
	}

	r.logger.With(log.Fields{
		"id": id,
	}).Debug("new user created")
	dto.ToDomain(user)
	return nil
}

func (r *UserRepository) GetByID(ctx context.Context, id string) (domain.User, error) {
	// TODO implement me
	panic("implement me")
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (domain.User, error) {
	// TODO implement me
	panic("implement me")
}

func (r *UserRepository) SelectAll(ctx context.Context) ([]domain.User, error) {
	// TODO implement me
	panic("implement me")
}

func (r *UserRepository) SelectByName(ctx context.Context, name string) ([]domain.User, error) {
	// TODO implement me
	panic("implement me")
}

func (r *UserRepository) UpdateByID(ctx context.Context, user *domain.User) error {
	// TODO implement me
	panic("implement me")
}

func (r *UserRepository) DeleteByID(ctx context.Context, id string) error {
	// TODO implement me
	panic("implement me")
}
