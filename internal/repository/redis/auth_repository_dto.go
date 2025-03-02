package redis

import (
	"time"

	"github.com/adanyl0v/pocket-ideas/internal/domain"
)

type saveSessionDto struct {
	ID           string             `json:"id"`
	UserID       string             `json:"user_id"`
	Fingerprint  domain.Fingerprint `json:"fingerprint"`
	RefreshToken string             `json:"refresh_token"`
	ExpiresAt    time.Time          `json:"expires_at"`
	CreatedAt    time.Time          `json:"created_at"`
	UpdatedAt    time.Time          `json:"updated_at"`
}

func newSaveSessionDto(s *domain.Session) saveSessionDto {
	return saveSessionDto{
		UserID:       s.User.ID,
		Fingerprint:  s.Fingerprint,
		RefreshToken: s.RefreshToken,
		ExpiresAt:    s.ExpiresAt,
	}
}

func (d *saveSessionDto) ToDomain(s *domain.Session) {
	s.ID = d.ID
	s.CreatedAt = d.CreatedAt
	s.UpdatedAt = d.UpdatedAt
}

type findSessionByIdDto saveSessionDto

func newFindSessionByIdDto(id string) findSessionByIdDto {
	return findSessionByIdDto{ID: id}
}

func (d *findSessionByIdDto) ToDomain(s *domain.Session) {
	s.User.ID = d.UserID
	s.Fingerprint = d.Fingerprint
	s.RefreshToken = d.RefreshToken
	s.ExpiresAt = d.ExpiresAt
	s.CreatedAt = d.CreatedAt
	s.UpdatedAt = d.UpdatedAt
}

type findSessionByRefreshTokenDto saveSessionDto

func newFindSessionByRefreshTokenDto(refreshToken string) findSessionByRefreshTokenDto {
	return findSessionByRefreshTokenDto{
		RefreshToken: refreshToken,
	}
}

func (d *findSessionByRefreshTokenDto) ToDomain(s *domain.Session) {
	s.ID = d.ID
	s.User.ID = d.UserID
	s.Fingerprint = d.Fingerprint
	s.ExpiresAt = d.ExpiresAt
	s.CreatedAt = d.CreatedAt
	s.UpdatedAt = d.UpdatedAt
}

type findSessionByFingerprintDto saveSessionDto

func newFindSessionByFingerprintDto(fp domain.Fingerprint) findSessionByFingerprintDto {
	return findSessionByFingerprintDto{
		Fingerprint: fp,
	}
}

func (d *findSessionByFingerprintDto) ToDomain(s *domain.Session) {
	s.ID = d.ID
	s.User.ID = d.UserID
	s.RefreshToken = d.RefreshToken
	s.ExpiresAt = d.ExpiresAt
	s.CreatedAt = d.CreatedAt
	s.UpdatedAt = d.UpdatedAt
}

type findAllSessionsDto saveSessionDto

func newFindAllSessionsDto() findAllSessionsDto {
	return findAllSessionsDto{}
}

func (d *findAllSessionsDto) ToDomain(s *domain.Session) {
	s.ID = d.ID
	s.User.ID = d.UserID
	s.Fingerprint = d.Fingerprint
	s.RefreshToken = d.RefreshToken
	s.ExpiresAt = d.ExpiresAt
	s.CreatedAt = d.CreatedAt
	s.UpdatedAt = d.UpdatedAt
}

type findSessionsByUserIdDto saveSessionDto

func newFindSessionsByUserIdDto(userId string) findSessionsByUserIdDto {
	return findSessionsByUserIdDto{
		UserID: userId,
	}
}

func (d *findSessionsByUserIdDto) ToDomain(s *domain.Session) {
	s.ID = d.ID
	s.Fingerprint = d.Fingerprint
	s.RefreshToken = d.RefreshToken
	s.ExpiresAt = d.ExpiresAt
	s.CreatedAt = d.CreatedAt
	s.UpdatedAt = d.UpdatedAt
}

type updateSessionByIdDto saveSessionDto

func newUpdateSessionByIdDto(s *domain.Session) updateSessionByIdDto {
	return updateSessionByIdDto{
		ID:           s.ID,
		UserID:       s.User.ID,
		Fingerprint:  s.Fingerprint,
		RefreshToken: s.RefreshToken,
		ExpiresAt:    s.ExpiresAt,
		CreatedAt:    s.CreatedAt,
	}
}

func (d *updateSessionByIdDto) ToDomain(s *domain.Session) {
	s.UpdatedAt = d.UpdatedAt
}
