package postgres

import (
	"github.com/adanyl0v/pocket-ideas/internal/domain"
	"github.com/jackc/pgx/v5/pgtype/zeronull"
	"time"
)

type saveUserDto struct {
	ID        string        `json:"id" db:"id"`
	Name      zeronull.Text `json:"name" db:"name"`
	Email     zeronull.Text `json:"email" db:"email"`
	Password  zeronull.Text `json:"password" db:"password"`
	CreatedAt time.Time     `json:"created_at" db:"created_at"`
	UpdatedAt time.Time     `json:"updated_at" db:"updated_at"`
}

func newSaveUserDto(u *domain.User) saveUserDto {
	return saveUserDto{
		Name:     zeronull.Text(u.Name),
		Email:    zeronull.Text(u.Email),
		Password: zeronull.Text(u.Password),
	}
}

func (d *saveUserDto) ToDomain(u *domain.User) {
	u.ID = d.ID
	u.CreatedAt = d.CreatedAt
	u.UpdatedAt = d.UpdatedAt
}
