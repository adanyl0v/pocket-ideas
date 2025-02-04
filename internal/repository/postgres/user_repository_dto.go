package postgres

import (
	"github.com/adanyl0v/pocket-ideas/internal/domain"
	"github.com/jackc/pgx/v5/pgtype"
	"time"
)

type userRepositoryDTO interface {
	ToDomain(user *domain.User)
	FromDomain(user *domain.User)
}

type createUserDTO struct {
	ID        string    `db:"id"`
	Name      string    `db:"name"`
	Email     string    `db:"email"`
	Password  string    `db:"password"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (dto *createUserDTO) ToDomain(user *domain.User) {
	user.ID = dto.ID
	user.CreatedAt = dto.CreatedAt
	user.UpdatedAt = dto.UpdatedAt
}

func (dto *createUserDTO) FromDomain(user *domain.User) {
	dto.Name = user.Name
	dto.Email = user.Email
	dto.Password = user.Password
}

type getUserByIdDTO struct {
	ID        string    `db:"id"`
	Name      string    `db:"name"`
	Email     string    `db:"email"`
	Password  string    `db:"password"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (dto *getUserByIdDTO) ToDomain(user *domain.User) {
	user.Name = dto.Name
	user.Email = dto.Email
	user.Password = dto.Password
	user.CreatedAt = dto.CreatedAt
	user.UpdatedAt = dto.UpdatedAt
}

func (dto *getUserByIdDTO) FromDomain(user *domain.User) {
	dto.ID = user.ID
}

type getUserByEmailDTO struct {
	ID        string    `db:"id"`
	Name      string    `db:"name"`
	Email     string    `db:"email"`
	Password  string    `db:"password"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (dto *getUserByEmailDTO) ToDomain(user *domain.User) {
	user.ID = dto.ID
	user.Name = dto.Name
	user.Password = dto.Password
	user.CreatedAt = dto.CreatedAt
	user.UpdatedAt = dto.UpdatedAt
}

func (dto *getUserByEmailDTO) FromDomain(user *domain.User) {
	dto.Email = user.Email
}

type selectAllUsersDTO struct {
	ID        string    `db:"id"`
	Name      string    `db:"name"`
	Email     string    `db:"email"`
	Password  string    `db:"password"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (dto *selectAllUsersDTO) ToDomain(user *domain.User) {
	user.ID = dto.ID
	user.Name = dto.Name
	user.Email = dto.Email
	user.Password = dto.Password
	user.CreatedAt = dto.CreatedAt
	user.UpdatedAt = dto.UpdatedAt
}

func (dto *selectAllUsersDTO) FromDomain(_ *domain.User) {}

type selectUsersByNameDTO struct {
	ID        string    `db:"id"`
	Name      string    `db:"name"`
	Email     string    `db:"email"`
	Password  string    `db:"password"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (dto *selectUsersByNameDTO) ToDomain(user *domain.User) {
	user.ID = dto.ID
	user.Email = dto.Email
	user.Password = dto.Password
	user.CreatedAt = dto.CreatedAt
	user.UpdatedAt = dto.UpdatedAt
}

func (dto *selectUsersByNameDTO) FromDomain(user *domain.User) {
	dto.Name = user.Name
}

type updateUserByIdDTO struct {
	ID        string      `db:"id"`
	Name      pgtype.Text `db:"name"`
	Email     pgtype.Text `db:"email"`
	Password  pgtype.Text `db:"password"`
	UpdatedAt time.Time   `db:"updated_at"`
}

func (dto *updateUserByIdDTO) ToDomain(user *domain.User) {
	user.UpdatedAt = dto.UpdatedAt
}

func (dto *updateUserByIdDTO) FromDomain(user *domain.User) {
	dto.ID = user.ID

	if user.Name != "" {
		dto.Name = pgtype.Text{String: user.Name, Valid: true}
	}
	if user.Email != "" {
		dto.Email = pgtype.Text{String: user.Email, Valid: true}
	}
	if user.Password != "" {
		dto.Password = pgtype.Text{String: user.Password, Valid: true}
	}
}
