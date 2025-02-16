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
	accessTokensWhiteListKey = "access_tokens_whitelist"
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

	if err := r.whiteListConn.Set(ctx, convertAccessTokenToWhiteListKey(accessToken),
		accessToken, expiration); err != nil {
		logger.WithError(err).Error("failed to save access token to white list")
		return err
	}

	logger.Debug("saved access token to whitelist")
	return nil
}

func (r *AuthRepository) FindAccessTokenInWhiteList(ctx context.Context, accessToken string) (bool, error) {
	logger := r.logger.With(log.Fields{"access_token": accessToken})

	n, err := r.whiteListConn.Exists(ctx, convertAccessTokenToWhiteListKey(accessToken))
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

	if err := r.whiteListConn.Delete(ctx, convertAccessTokenToWhiteListKey(accessToken)); err != nil {
		logger.WithError(err).Error("failed to delete access token from white list")
		return err
	}

	logger.Debug("deleted access token from white list")
	return nil
}

func (r *AuthRepository) SaveRefreshTokenToBlackList(ctx context.Context, refreshToken string, expiration time.Duration) error {
	// TODO implement me
	panic("implement me")
}

func (r *AuthRepository) FindRefreshTokenInBlackList(ctx context.Context, refreshToken string) (bool, error) {
	// TODO implement me
	panic("implement me")
}

func (r *AuthRepository) DeleteRefreshTokenFromBlackList(ctx context.Context, refreshToken string) error {
	// TODO implement me
	panic("implement me")
}

func convertAccessTokenToWhiteListKey(accessToken string) string {
	return fmt.Sprint(accessTokensWhiteListKey, ":", accessToken)
}
