package postgres

import (
	"context"
	"errors"
	"github.com/adanyl0v/pocket-ideas/internal/domain"
	"github.com/adanyl0v/pocket-ideas/internal/repository"
	pgdb "github.com/adanyl0v/pocket-ideas/pkg/database/postgres"
	"github.com/adanyl0v/pocket-ideas/pkg/log"
	"github.com/adanyl0v/pocket-ideas/pkg/proxerr"
	uuidgen "github.com/adanyl0v/pocket-ideas/pkg/uuid"
	"time"
)

var (
	ErrUserExists   = errors.New(msgUserAlreadyExists)
	ErrUserNotFound = errors.New(msgUserNotFound)
	ErrNoUsersFound = errors.New(msgNoUsersFound)

	errGenerateUserId = errors.New(msgFailedToGenerateUserID)
)

// For testing purpose
const (
	msgUserAlreadyExists      = "user already exists"
	msgUserNotFound           = "user not found"
	msgNoUsersFound           = "no users found"
	msgFailedToGenerateUserID = "failed to generate user id"
	msgFailedToSelectAllUsers = "failed to select all users"
	msgFailedToScanAllUsers   = "failed to scan all users"
	msgNoUserRowsAffected     = "no rows affected"
)

type UserRepository struct {
	Repository
	idGen uuidgen.Generator
}

func NewUserRepository(conn pgdb.Conn, logger log.Logger, idGen uuidgen.Generator) *UserRepository {
	return &UserRepository{
		Repository: Repository{
			conn:   conn,
			logger: logger,
		},
		idGen: idGen,
	}
}

func (r *UserRepository) Begin(ctx context.Context) (repository.Tx, error) {
	tx, err := r.conn.Begin(ctx)
	if err != nil {
		r.logger.WithError(err).Error("failed to begin transaction")
		return nil, err
	}

	return tx, nil
}

func (r *UserRepository) WithTx(tx repository.Tx) repository.Repository {
	conn, ok := tx.(pgdb.Tx)
	if !ok {
		r.logger.With(log.Fields{"transaction": tx}).Warn("unsupported transaction type")
		return nil
	}

	return &UserRepository{
		Repository: Repository{
			conn:   conn,
			logger: r.logger,
		},
		idGen: r.idGen,
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
		err = proxerr.New(errGenerateUserId, err.Error())
		r.logger.WithError(err).Error(errGenerateUserId.Error())
		return err
	}

	dto.ID = id
	dto.CreatedAt = time.Now()
	dto.UpdatedAt = dto.CreatedAt
	if err = r.conn.Exec(ctx, insertUserQuery, dto.ID, dto.Name, dto.Email,
		dto.Password, dto.CreatedAt, dto.UpdatedAt); err != nil {

		var pxErr proxerr.Error
		if errors.As(err, &pxErr) {
			if errors.Is(pxErr.Unwrap(), pgdb.ErrUniqueViolation) {
				r.logger.With(log.Fields{
					"email": dto.Email,
				}).Warn("tried to create a new user with existing email")
				return proxerr.New(ErrUserExists, pxErr.Error())
			}
		}

		r.logger.WithError(err).Error("failed to insert user")
		return err
	}

	r.logger.With(log.Fields{
		"id": id,
	}).Debug("new user created")
	dto.ToDomain(user)
	return nil
}

const getUserByIdQuery = `
SELECT name, email, password, created_at, updated_at
FROM users WHERE id = $1
`

func (r *UserRepository) GetByID(ctx context.Context, id string) (domain.User, error) {
	var user = domain.User{ID: id}
	var dto getUserByIdDTO
	dto.FromDomain(&user)

	if err := r.conn.QueryRow(ctx, getUserByIdQuery, dto.ID).Scan(&dto.Name, &dto.Email,
		&dto.Password, &dto.CreatedAt, &dto.UpdatedAt); err != nil {

		var pxErr proxerr.Error
		if errors.As(err, &pxErr) {
			if errors.Is(pxErr.Unwrap(), pgdb.ErrNoRows) {
				r.logger.With(log.Fields{
					"id": id,
				}).Warn("tried to get a nonexistent user")
				return domain.User{}, proxerr.New(ErrUserNotFound, pxErr.Error())
			}
		}

		r.logger.WithError(err).Error("failed to get the user by id")
		return domain.User{}, err
	}

	r.logger.With(log.Fields{
		"id": id,
	}).Debug("found user by id")
	dto.ToDomain(&user)
	return user, nil
}

const getUserByEmailQuery = `
SELECT id, name, password, created_at, updated_at
FROM users WHERE email = $1
`

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (domain.User, error) {
	var user = domain.User{Email: email}
	var dto getUserByEmailDTO
	dto.FromDomain(&user)

	if err := r.conn.QueryRow(ctx, getUserByEmailQuery, dto.Email).Scan(&dto.ID, &dto.Name,
		&dto.Password, &dto.CreatedAt, &dto.UpdatedAt); err != nil {

		var pxErr proxerr.Error
		if errors.As(err, &pxErr) {
			if errors.Is(pxErr.Unwrap(), pgdb.ErrNoRows) {
				r.logger.With(log.Fields{
					"email": email,
				}).Warn("tried to get a nonexistent user")
				return domain.User{}, proxerr.New(ErrUserNotFound, pxErr.Error())
			}
		}

		r.logger.WithError(err).Error("failed to get the user by email")
		return domain.User{}, err
	}

	r.logger.With(log.Fields{
		"id":    dto.ID,
		"email": email,
	}).Debug("found user by email")
	dto.ToDomain(&user)
	return user, nil
}

const selectAllUsersQuery = `
SELECT id, name, email, password, created_at, updated_at
FROM users
`

func (r *UserRepository) SelectAll(ctx context.Context) ([]domain.User, error) {
	rows, err := r.conn.Query(ctx, selectAllUsersQuery)
	if err != nil {
		r.logger.WithError(err).Error(msgFailedToSelectAllUsers)
		return nil, err
	}
	defer rows.Close()

	next := rows.Next()
	if !next {
		err = rows.Err()
		if err == nil {
			err = proxerr.New(ErrNoUsersFound, pgdb.ErrNoRows.Error())
			r.logger.Warn(ErrNoUsersFound.Error())
		} else {
			r.logger.WithError(err).Error(msgNoUserRowsAffected)
		}

		return nil, err
	}

	users := make([]domain.User, 0, 4)
	for next {
		var dto selectAllUsersDTO
		if err = rows.Scan(&dto.ID, &dto.Name, &dto.Email, &dto.Password,
			&dto.CreatedAt, &dto.UpdatedAt); err != nil {

			r.logger.WithError(err).Error(msgFailedToScanAllUsers)
			return nil, err
		}

		var user domain.User
		dto.ToDomain(&user)
		users = append(users, user)

		next = rows.Next()
	}

	r.logger.With(log.Fields{"amount": len(users)}).Debug("got all users")
	return users, nil
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
