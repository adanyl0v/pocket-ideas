package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/adanyl0v/pocket-ideas/pkg/proxerr"
	"strings"
	"time"

	"github.com/adanyl0v/pocket-ideas/internal/domain"
	"github.com/adanyl0v/pocket-ideas/pkg/cache"
	"github.com/adanyl0v/pocket-ideas/pkg/log"
	"github.com/adanyl0v/pocket-ideas/pkg/uuid"
)

var ErrNotFound = errors.New("not found")

type JSONer interface {
	Marshal(v interface{}) ([]byte, error)
	Unmarshal(data []byte, v interface{}) error
}

type stdJSONer struct{}

func (j *stdJSONer) Marshal(v interface{}) ([]byte, error)      { return json.Marshal(v) }
func (j *stdJSONer) Unmarshal(data []byte, v interface{}) error { return json.Unmarshal(data, v) }

const (
	// sessionKeyFormat will be interpreted as "session:<session_id>"
	sessionKeyFormat = "session:%s"

	// whitelistKeyFormat will be interpreted as "whitelist:<jwt_token>"
	whitelistKeyFormat = "whitelist:%s"

	// blacklistKeyFormat will be interpreted as "blacklist:<jwt_token>"
	blacklistKeyFormat = "blacklist:%s"
)

type AuthRepository struct {
	sessionsConn  cache.Conn
	whitelistConn cache.Conn
	blacklistConn cache.Conn
	logger        log.Logger
	idGen         uuid.Generator
	jsoner        JSONer
}

func NewAuthRepository(
	sessionsConn cache.Conn,
	whitelistConn cache.Conn,
	blacklistConn cache.Conn,
	logger log.Logger,
	idGen uuid.Generator,
) *AuthRepository {
	return &AuthRepository{
		sessionsConn:  sessionsConn,
		whitelistConn: whitelistConn,
		blacklistConn: blacklistConn,
		logger:        logger,
		idGen:         idGen,
		jsoner:        new(stdJSONer),
	}
}

func (r *AuthRepository) SetJSONer(jsoner JSONer) {
	r.jsoner = jsoner
}

func (r *AuthRepository) SaveSession(ctx context.Context, session *domain.Session) error {
	dto := newSaveSessionDto(session)

	var err error
	dto.ID, err = r.idGen.NewV7()
	if err != nil {
		r.logger.WithError(err).Error("failed to generate session uuid")
		return err
	}

	// Set the UTC timezone explicitly
	dto.CreatedAt = time.Now().UTC()
	dto.UpdatedAt = dto.CreatedAt

	b, err := r.jsoner.Marshal(dto)
	if err != nil {
		r.logger.WithError(err).Error("failed to marshal a dto")
		return err
	}

	if err = r.sessionsConn.Set(ctx, formatToSessionKey(dto.ID), b, -1); err != nil {
		r.logger.WithError(err).Error("failed to save a session")
		return err
	}

	dto.ToDomain(session)
	r.logger.With(log.Fields{"id": dto.ID}).Debug("saved a session")
	return nil
}

func (r *AuthRepository) FindSessionById(ctx context.Context, id string) (domain.Session, error) {
	logger := r.logger.With(log.Fields{"id": id})

	var session domain.Session
	dto := newFindSessionByIdDto(id)

	var raw string
	if err := r.sessionsConn.Get(ctx, formatToSessionKey(dto.ID), &raw); err != nil {
		logger.WithError(err).Error("failed to find a session by id")
		return domain.Session{}, proxerr.New(ErrNotFound, err.Error())
	}

	if err := r.jsoner.Unmarshal([]byte(raw), &dto); err != nil {
		logger.WithError(err).Error("failed to unmarshal a session")
		return domain.Session{}, err
	}

	dto.ToDomain(&session)
	r.logger.With(log.Fields{
		"id":      dto.ID,
		"user_id": dto.UserID,
	}).Debug("found the session by id")
	return session, nil
}

func (r *AuthRepository) FindSessionByRefreshToken(ctx context.Context, refreshToken string) (domain.Session, error) {
	logger := r.logger.With(log.Fields{"refresh_token": refreshToken})

	var session domain.Session
	dto := newFindSessionByRefreshTokenDto(refreshToken)

	it := r.sessionsConn.Scan(ctx, cache.DefaultScanner)
	if err := it.Err(); err != nil {
		logger.WithError(it.Err()).Error("failed to scan sessions")
		return domain.Session{}, err
	}

	var raw string
	for it.Next(ctx) {
		var s string
		if err := r.sessionsConn.Get(ctx, it.Val(), &s); err != nil {
			logger.WithError(err).Error("failed to get a session by key")
			return domain.Session{}, err
		}

		// Trying to avoid json unmarshalling on every iteration
		key := fmt.Sprintf("\"refresh_token\":\"%s\"", refreshToken)
		if strings.Contains(s, key) {
			raw = s
			break
		}
	}
	// All tests will return an error here, because the raw string is always empty
	if raw == "" {
		logger.WithError(ErrNotFound).Error("failed to find a session by refresh token")
		return domain.Session{}, ErrNotFound
	}

	if err := r.jsoner.Unmarshal([]byte(raw), &dto); err != nil {
		logger.WithError(err).Error("failed to unmarshal a session")
		return domain.Session{}, err
	}

	dto.ToDomain(&session)
	r.logger.With(log.Fields{
		"id":      dto.ID,
		"user_id": dto.UserID,
	}).Debug("found the session by refresh token")
	return session, nil
}

func (r *AuthRepository) FindSessionByFingerprint(ctx context.Context, fp domain.Fingerprint) (domain.Session, error) {
	logger := r.logger.With(log.Fields{"fingerprint": fp})

	var session domain.Session
	dto := newFindSessionByFingerprintDto(fp)

	fpBytes, err := r.jsoner.Marshal(fp)
	if err != nil {
		logger.WithError(err).Error("failed to marshal a fingerprint")
		return domain.Session{}, err
	}
	fpRaw := string(fpBytes)

	it := r.sessionsConn.Scan(ctx, cache.DefaultScanner)
	if err = it.Err(); err != nil {
		logger.WithError(it.Err()).Error("failed to scan sessions")
		return domain.Session{}, err
	}

	var raw string
	for it.Next(ctx) {
		var s string
		if err = r.sessionsConn.Get(ctx, it.Val(), &s); err != nil {
			logger.WithError(err).Error("failed to get a session by key")
			return domain.Session{}, err
		}

		// Trying to avoid json unmarshalling on every iteration
		if strings.Contains(s, fpRaw) {
			raw = s
			break
		}
	}
	// All tests will return an error here, because the raw string is always empty
	if raw == "" {
		logger.WithError(ErrNotFound).Error("failed to find a session by fingerprint")
		return domain.Session{}, ErrNotFound
	}

	if err = r.jsoner.Unmarshal([]byte(raw), &dto); err != nil {
		logger.WithError(err).Error("failed to unmarshal a session")
		return domain.Session{}, err
	}

	dto.ToDomain(&session)
	r.logger.With(log.Fields{
		"id":      dto.ID,
		"user_id": dto.UserID,
	}).Debug("found the session by fingerprint")
	return session, nil
}

// FindAllSessions returns a zero-length slice if no sessions were found
func (r *AuthRepository) FindAllSessions(ctx context.Context) ([]domain.Session, error) {
	it := r.sessionsConn.Scan(ctx, cache.DefaultScanner)
	if err := it.Err(); err != nil {
		r.logger.WithError(err).Error("failed to scan sessions")
		return nil, err
	}

	sessions := make([]domain.Session, 0)
	for it.Next(ctx) {
		var session domain.Session
		dto := newFindAllSessionsDto()

		var raw string
		if err := r.sessionsConn.Get(ctx, it.Val(), &raw); err != nil {
			r.logger.WithError(err).Error("failed to get a session by key")
			return nil, err
		}

		if err := r.jsoner.Unmarshal([]byte(raw), &dto); err != nil {
			r.logger.WithError(err).Error("failed to unmarshal a session")
			return nil, err
		}

		dto.ToDomain(&session)
		sessions = append(sessions, session)
	}

	r.logger.Debug(fmt.Sprintf("found %d sessions", len(sessions)))
	return sessions, nil
}

// FindSessionsByUserId returns a zero-length slice if no sessions were found
func (r *AuthRepository) FindSessionsByUserId(ctx context.Context, userId string) ([]domain.Session, error) {
	logger := r.logger.With(log.Fields{"user_id": userId})

	it := r.sessionsConn.Scan(ctx, cache.DefaultScanner)
	if err := it.Err(); err != nil {
		logger.WithError(err).Error("failed to scan sessions")
		return nil, err
	}

	sessions := make([]domain.Session, 0)
	key := fmt.Sprintf("\"user_id\":\"%s\"", userId)
	for it.Next(ctx) {
		var raw string
		if err := r.sessionsConn.Get(ctx, it.Val(), &raw); err != nil {
			logger.WithError(err).Error("failed to get a session by key")
			return nil, err
		}

		if strings.Contains(raw, key) {
			var session domain.Session
			dto := newFindSessionsByUserIdDto(userId)

			if err := r.jsoner.Unmarshal([]byte(raw), &dto); err != nil {
				logger.WithError(err).Error("failed to unmarshal a session")
				return nil, err
			}

			dto.ToDomain(&session)
			sessions = append(sessions, session)
		}
	}

	logger.Debug(fmt.Sprintf("found %d sessions by user id", len(sessions)))
	return sessions, nil
}

func (r *AuthRepository) UpdateSessionById(ctx context.Context, session *domain.Session) error {
	logger := r.logger.With(log.Fields{"id": session.ID})

	dto := newUpdateSessionByIdDto(session)
	dto.UpdatedAt = time.Now().UTC()

	b, err := r.jsoner.Marshal(dto)
	if err != nil {
		logger.WithError(err).Error("failed to marshal a session")
		return err
	}

	if err = r.sessionsConn.Set(ctx, formatToSessionKey(dto.ID), b, -1); err != nil {
		logger.WithError(err).Error("failed to update a session")
		return err
	}

	dto.ToDomain(session)
	logger.Debug("updated a session")
	return nil
}

func (r *AuthRepository) DeleteSessionById(ctx context.Context, id string) error {
	logger := r.logger.With(log.Fields{"id": id})

	_, err := r.sessionsConn.Delete(ctx, formatToSessionKey(id))
	if err != nil {
		logger.WithError(err).Error("failed to delete a session")
		return err
	}

	logger.Debug("deleted a session")
	return nil
}

func (r *AuthRepository) SaveAccessTokenToWhitelist(ctx context.Context, accessToken string, expiration time.Duration) error {
	logger := r.logger.With(log.Fields{"access_token": accessToken})

	if err := r.whitelistConn.Set(ctx, formatAccessTokenIntoCacheKey(accessToken),
		accessToken, expiration); err != nil {
		logger.WithError(err).Error("failed to save access token to whitelist")
		return err
	}

	logger.Debug("saved access token to whitelist")
	return nil
}

func (r *AuthRepository) FindAccessTokenInWhitelist(ctx context.Context, accessToken string) (bool, error) {
	logger := r.logger.With(log.Fields{"access_token": accessToken})

	n, err := r.whitelistConn.Exists(ctx, formatAccessTokenIntoCacheKey(accessToken))
	if err != nil {
		logger.WithError(err).Error("failed to find access token in whitelist")
		return false, err
	}

	if n < 1 {
		logger.Debug("access token not found in whitelist")
		return false, nil
	}

	logger.Debug("found access token in whitelist")
	return true, nil
}

func (r *AuthRepository) DeleteAccessTokenFromWhitelist(ctx context.Context, accessToken string) error {
	logger := r.logger.With(log.Fields{"access_token": accessToken})

	if _, err := r.whitelistConn.Delete(ctx, formatAccessTokenIntoCacheKey(accessToken)); err != nil {
		logger.WithError(err).Error("failed to delete access token from whitelist")
		return err
	}

	logger.Debug("deleted access token from whitelist")
	return nil
}

func (r *AuthRepository) SaveRefreshTokenToBlacklist(ctx context.Context, refreshToken string, expiration time.Duration) error {
	logger := r.logger.With(log.Fields{"refresh_token": refreshToken})

	if err := r.blacklistConn.Set(ctx, formatRefreshTokenIntoCacheKey(refreshToken),
		refreshToken, expiration); err != nil {
		logger.WithError(err).Error("failed to save refresh token to blacklist")
		return err
	}

	logger.Debug("saved refresh token to blacklist")
	return nil
}

func (r *AuthRepository) FindRefreshTokenInBlacklist(ctx context.Context, refreshToken string) (bool, error) {
	logger := r.logger.With(log.Fields{"refresh_token": refreshToken})

	n, err := r.blacklistConn.Exists(ctx, formatRefreshTokenIntoCacheKey(refreshToken))
	if err != nil {
		logger.WithError(err).Error("failed to find refresh token in blacklist")
		return false, err
	}

	if n < 1 {
		logger.Debug("refresh token not found in blacklist")
		return false, nil
	}

	logger.Debug("found refresh token in blacklist")
	return true, nil
}

func (r *AuthRepository) DeleteRefreshTokenFromBlacklist(ctx context.Context, refreshToken string) error {
	logger := r.logger.With(log.Fields{"refresh_token": refreshToken})

	if _, err := r.blacklistConn.Delete(ctx, formatRefreshTokenIntoCacheKey(refreshToken)); err != nil {
		logger.WithError(err).Error("failed to delete refresh token from blacklist")
		return err
	}

	logger.Debug("deleted refresh token from blacklist")
	return nil
}

func formatToSessionKey(sessionId string) string {
	if sessionId == "" {
		sessionId = "*"
	}

	return fmt.Sprintf(sessionKeyFormat, sessionId)
}

func formatAccessTokenIntoCacheKey(accessToken string) string {
	return fmt.Sprintf(whitelistKeyFormat, accessToken)
}

func formatRefreshTokenIntoCacheKey(refreshToken string) string {
	return fmt.Sprintf(blacklistKeyFormat, refreshToken)
}
