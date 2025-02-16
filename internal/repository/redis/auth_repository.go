package redis

import (
	"context"
	"fmt"
	"github.com/adanyl0v/pocket-ideas/internal/domain"
	"github.com/adanyl0v/pocket-ideas/pkg/cache"
	"github.com/adanyl0v/pocket-ideas/pkg/log"
	"github.com/adanyl0v/pocket-ideas/pkg/uuid"
	"time"
)

const (
	// whitelistKeyFormat will be interpreted as "whitelist:<jwt_token>"
	whitelistKeyFormat = "whitelist:%s"

	// blacklistKeyFormat will be interpreted as "blacklist:<jwt_token>"
	blacklistKeyFormat = "blacklist:%s"
)

type AuthRepository struct {
	sessionsConn  cache.Conn
	whiteListConn cache.Conn
	blackListConn cache.Conn
	logger        log.Logger
	idGen         uuid.Generator
}

func NewAuthRepository(
	sessionsConn cache.Conn,
	whiteListConn cache.Conn,
	blackListConn cache.Conn,
	logger log.Logger,
	idGen uuid.Generator) *AuthRepository {
	return &AuthRepository{
		sessionsConn:  sessionsConn,
		whiteListConn: whiteListConn,
		blackListConn: blackListConn,
		logger:        logger,
		idGen:         idGen,
	}
}

func (r *AuthRepository) SaveSession(ctx context.Context, session *domain.Session) error {
	// TODO implement me
	panic("implement me")
}

func (r *AuthRepository) FindSessionById(ctx context.Context, id string) (domain.Session, error) {
	// TODO implement me
	panic("implement me")
}

func (r *AuthRepository) FindSessionByRefreshToken(ctx context.Context, refreshToken string) (domain.Session, error) {
	// TODO implement me
	panic("implement me")
}

func (r *AuthRepository) FindSessionByFingerprint(ctx context.Context, fp domain.Fingerprint) (domain.Session, error) {
	// TODO implement me
	panic("implement me")
}

func (r *AuthRepository) FindAllSessions(ctx context.Context) ([]domain.Session, error) {
	// TODO implement me
	panic("implement me")
}

func (r *AuthRepository) FindSessionByUserId(ctx context.Context, userId string) ([]domain.Session, error) {
	// TODO implement me
	panic("implement me")
}

func (r *AuthRepository) UpdateSessionById(ctx context.Context, session *domain.Session) error {
	// TODO implement me
	panic("implement me")
}

func (r *AuthRepository) DeleteSessionById(ctx context.Context, id string) error {
	// TODO implement me
	panic("implement me")
}

func (r *AuthRepository) SaveAccessTokenToWhiteList(ctx context.Context, accessToken string, expiration time.Duration) error {
	logger := r.logger.With(log.Fields{"access_token": accessToken})

	if err := r.whiteListConn.Set(ctx, formatAccessTokenIntoCacheKey(accessToken),
		accessToken, expiration); err != nil {
		logger.WithError(err).Error("failed to save access token to whitelist")
		return err
	}

	logger.Debug("saved access token to whitelist")
	return nil
}

func (r *AuthRepository) FindAccessTokenInWhiteList(ctx context.Context, accessToken string) (bool, error) {
	logger := r.logger.With(log.Fields{"access_token": accessToken})

	n, err := r.whiteListConn.Exists(ctx, formatAccessTokenIntoCacheKey(accessToken))
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

func (r *AuthRepository) DeleteAccessTokenFromWhiteList(ctx context.Context, accessToken string) error {
	logger := r.logger.With(log.Fields{"access_token": accessToken})

	if _, err := r.whiteListConn.Delete(ctx, formatAccessTokenIntoCacheKey(accessToken)); err != nil {
		logger.WithError(err).Error("failed to delete access token from whitelist")
		return err
	}

	logger.Debug("deleted access token from whitelist")
	return nil
}

func (r *AuthRepository) SaveRefreshTokenToBlackList(ctx context.Context, refreshToken string, expiration time.Duration) error {
	logger := r.logger.With(log.Fields{"refresh_token": refreshToken})

	if err := r.blackListConn.Set(ctx, formatRefreshTokenIntoCacheKey(refreshToken),
		refreshToken, expiration); err != nil {
		logger.WithError(err).Error("failed to save refresh token to blacklist")
		return err
	}

	logger.Debug("saved refresh token to blacklist")
	return nil
}

func (r *AuthRepository) FindRefreshTokenInBlackList(ctx context.Context, refreshToken string) (bool, error) {
	logger := r.logger.With(log.Fields{"refresh_token": refreshToken})

	n, err := r.blackListConn.Exists(ctx, formatRefreshTokenIntoCacheKey(refreshToken))
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

func (r *AuthRepository) DeleteRefreshTokenFromBlackList(ctx context.Context, refreshToken string) error {
	logger := r.logger.With(log.Fields{"refresh_token": refreshToken})

	if _, err := r.blackListConn.Delete(ctx, formatRefreshTokenIntoCacheKey(refreshToken)); err != nil {
		logger.WithError(err).Error("failed to delete refresh token from blacklist")
		return err
	}

	logger.Debug("deleted refresh token from blacklist")
	return nil
}

func formatAccessTokenIntoCacheKey(accessToken string) string {
	return fmt.Sprintf(whitelistKeyFormat, accessToken)
}

func formatRefreshTokenIntoCacheKey(refreshToken string) string {
	return fmt.Sprintf(blacklistKeyFormat, refreshToken)
}
