package postgres

import (
	"github.com/adanyl0v/pocket-ideas/internal/domain"
	"time"
)

type userDTO interface {
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
