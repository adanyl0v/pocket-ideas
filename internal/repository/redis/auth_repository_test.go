package redis

import (
	"context"
	"errors"
	"github.com/adanyl0v/pocket-ideas/internal/domain"
	_redisrepoMock "github.com/adanyl0v/pocket-ideas/internal/repository/redis/mocks"
	_cacheMock "github.com/adanyl0v/pocket-ideas/mocks/pkg/cache"
	_uuidMock "github.com/adanyl0v/pocket-ideas/mocks/pkg/uuid"
	"github.com/adanyl0v/pocket-ideas/pkg/log/slog"
	"github.com/golang/mock/gomock"
	slogzap "github.com/samber/slog-zap/v2"
	"github.com/stretchr/testify/require"
	stdslog "log/slog"
	"testing"
	"time"
)

type (
	authRepositoryTestCase struct {
		reg authRepositoryTestCaseRegister
		cmd authRepositoryTestCaseCommand
		exp authRepositoryTestCaseExpect
	}

	authRepositoryTestCaseRegister func(_ *gomock.Controller, conn *_cacheMock.MockConn, _ *_uuidMock.MockGenerator)
	authRepositoryTestCaseCommand  func(repo *AuthRepository) error
	authRepositoryTestCaseExpect   func(err error)
)

type (
	authRepoSessionsTestCase struct {
		reg authRepoSessionsTcRegister
		cmd authRepoSessionsTcCommand
		exp authRepoSessionsTcExpect
	}

	authRepoSessionsTcRegister func(
		_ *gomock.Controller,
		conn *_cacheMock.MockConn,
		_ *_uuidMock.MockGenerator,
		jsoner *_redisrepoMock.MockJSONer,
	)
	authRepoSessionsTcCommand func(repo *AuthRepository) error
	authRepoSessionsTcExpect  func(err error)
)

func TestAuthRepository_SaveSession(t *testing.T) {
	tcs := map[string]authRepoSessionsTestCase{
		"SUCCESS": {
			reg: func(_ *gomock.Controller, conn *_cacheMock.MockConn, idGen *_uuidMock.MockGenerator, jsoner *_redisrepoMock.MockJSONer) {
				idGen.EXPECT().NewV7().Return("", nil)
				jsoner.EXPECT().Marshal(gomock.Any()).Return(nil, nil)
				conn.EXPECT().Set(gomock.Any(), formatToSessionKey(""), gomock.Any(), gomock.Any()).
					Return(nil)
			},
			cmd: func(repo *AuthRepository) error {
				return repo.SaveSession(context.Background(), new(domain.Session))
			},
			exp: func(err error) {
				require.NoError(t, err)
			},
		},
		"FAILED to generate session uuid": {
			reg: func(_ *gomock.Controller, _ *_cacheMock.MockConn, idGen *_uuidMock.MockGenerator, _ *_redisrepoMock.MockJSONer) {
				idGen.EXPECT().NewV7().Return("", errors.New(""))
			},
			cmd: func(repo *AuthRepository) error {
				return repo.SaveSession(context.Background(), new(domain.Session))
			},
			exp: func(err error) {
				require.Error(t, err)
			},
		},
		"FAILED to marshal a session": {
			reg: func(_ *gomock.Controller, _ *_cacheMock.MockConn, idGen *_uuidMock.MockGenerator, jsoner *_redisrepoMock.MockJSONer) {
				idGen.EXPECT().NewV7().Return("", nil)
				jsoner.EXPECT().Marshal(gomock.Any()).Return(nil, errors.New(""))
			},
			cmd: func(repo *AuthRepository) error {
				return repo.SaveSession(context.Background(), new(domain.Session))
			},
			exp: func(err error) {
				require.Error(t, err)
			},
		},
		"FAILED to save a session": {
			reg: func(_ *gomock.Controller, conn *_cacheMock.MockConn, idGen *_uuidMock.MockGenerator, jsoner *_redisrepoMock.MockJSONer) {
				idGen.EXPECT().NewV7().Return("", nil)
				jsoner.EXPECT().Marshal(gomock.Any()).Return(nil, nil)
				conn.EXPECT().Set(gomock.Any(), formatToSessionKey(""), gomock.Any(), gomock.Any()).
					Return(errors.New(""))
			},
			cmd: func(repo *AuthRepository) error {
				return repo.SaveSession(context.Background(), new(domain.Session))
			},
			exp: func(err error) {
				require.Error(t, err)
			},
		},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			runSessionsTestCase(t, &tc)
		})
	}
}

func TestAuthRepository_FindSessionById(t *testing.T) {
	tcs := map[string]authRepoSessionsTestCase{
		"SUCCESS": {
			reg: func(_ *gomock.Controller, conn *_cacheMock.MockConn, _ *_uuidMock.MockGenerator, jsoner *_redisrepoMock.MockJSONer) {
				jsoner.EXPECT().Unmarshal(gomock.Any(), gomock.Any()).Return(nil)
				conn.EXPECT().Get(gomock.Any(), formatToSessionKey(""), gomock.Any()).
					Return(nil)
			},
			cmd: func(repo *AuthRepository) error {
				_, err := repo.FindSessionById(context.Background(), "")
				return err
			},
			exp: func(err error) {
				require.NoError(t, err)
			},
		},
		"FAILED to find a session by id": {
			reg: func(_ *gomock.Controller, conn *_cacheMock.MockConn, _ *_uuidMock.MockGenerator, _ *_redisrepoMock.MockJSONer) {
				conn.EXPECT().Get(gomock.Any(), formatToSessionKey(""), gomock.Any()).
					Return(errors.New(""))
			},
			cmd: func(repo *AuthRepository) error {
				_, err := repo.FindSessionById(context.Background(), "")
				return err
			},
			exp: func(err error) {
				require.Error(t, err)
			},
		},
		"FAILED to unmarshal a session": {
			reg: func(_ *gomock.Controller, conn *_cacheMock.MockConn, _ *_uuidMock.MockGenerator, jsoner *_redisrepoMock.MockJSONer) {
				conn.EXPECT().Get(gomock.Any(), formatToSessionKey(""), gomock.Any()).Return(nil)
				jsoner.EXPECT().Unmarshal(gomock.Any(), gomock.Any()).Return(errors.New(""))
			},
			cmd: func(repo *AuthRepository) error {
				_, err := repo.FindSessionById(context.Background(), "")
				return err
			},
			exp: func(err error) {
				require.Error(t, err)
			},
		},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			runSessionsTestCase(t, &tc)
		})
	}
}

func runSessionsTestCase(t *testing.T, tc *authRepoSessionsTestCase) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	sessionsConn := _cacheMock.NewMockConn(ctrl)
	idGen := _uuidMock.NewMockGenerator(ctrl)
	jsoner := _redisrepoMock.NewMockJSONer(ctrl)
	if tc.reg != nil {
		tc.reg(ctrl, sessionsConn, idGen, jsoner)
	}

	logger := slog.NewLogger(stdslog.New(slogzap.Option{}.NewZapHandler()))
	repo := NewAuthRepository(sessionsConn, nil, nil, logger, idGen)
	repo.SetJSONer(jsoner)

	var err error
	if tc.cmd != nil {
		err = tc.cmd(repo)
	}

	if tc.exp != nil {
		tc.exp(err)
	}
}

func TestAuthRepository_SaveAccessTokenToWhiteList(t *testing.T) {
	const token = "testAccessToken"
	const expiration = time.Duration(-1)

	tcs := map[string]authRepositoryTestCase{
		"SUCCESS": {
			reg: func(_ *gomock.Controller, conn *_cacheMock.MockConn, _ *_uuidMock.MockGenerator) {
				conn.EXPECT().Set(gomock.Any(), formatAccessTokenIntoCacheKey(token), token, expiration).Return(nil)
			},
			cmd: func(repo *AuthRepository) error {
				return repo.SaveAccessTokenToWhitelist(context.Background(), token, expiration)
			},
			exp: func(err error) {
				require.NoError(t, err)
			},
		},
		"FAILURE": {
			reg: func(_ *gomock.Controller, conn *_cacheMock.MockConn, _ *_uuidMock.MockGenerator) {
				conn.EXPECT().Set(gomock.Any(), formatAccessTokenIntoCacheKey(token), token, expiration).Return(errors.New(""))
			},
			cmd: func(repo *AuthRepository) error {
				return repo.SaveAccessTokenToWhitelist(context.Background(), token, expiration)
			},
			exp: func(err error) {
				require.Error(t, err)
			},
		},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			runAccessTokensWhiteListTestCase(t, &tc)
		})
	}
}

func TestAuthRepository_FindAccessTokenInWhiteList(t *testing.T) {
	const token = "testAccessToken"

	tcs := map[string]authRepositoryTestCase{
		"SUCCESS": {
			reg: func(_ *gomock.Controller, conn *_cacheMock.MockConn, _ *_uuidMock.MockGenerator) {
				conn.EXPECT().Exists(gomock.Any(), formatAccessTokenIntoCacheKey(token)).
					Return(int64(1), nil)
			},
			cmd: func(repo *AuthRepository) error {
				found, err := repo.FindAccessTokenInWhitelist(context.Background(), token)
				require.True(t, found)
				return err
			},
			exp: func(err error) {
				require.NoError(t, err)
			},
		},
		"FAILURE": {
			reg: func(_ *gomock.Controller, conn *_cacheMock.MockConn, _ *_uuidMock.MockGenerator) {
				conn.EXPECT().Exists(gomock.Any(), formatAccessTokenIntoCacheKey(token)).
					Return(int64(0), errors.New(""))
			},
			cmd: func(repo *AuthRepository) error {
				found, err := repo.FindAccessTokenInWhitelist(context.Background(), token)
				require.False(t, found)
				return err
			},
			exp: func(err error) {
				require.Error(t, err)
			},
		},
		"FAILURE not found": {
			reg: func(_ *gomock.Controller, conn *_cacheMock.MockConn, _ *_uuidMock.MockGenerator) {
				conn.EXPECT().Exists(gomock.Any(), formatAccessTokenIntoCacheKey(token)).
					Return(int64(0), nil)
			},
			cmd: func(repo *AuthRepository) error {
				found, err := repo.FindAccessTokenInWhitelist(context.Background(), token)
				require.False(t, found)
				return err
			},
			exp: func(err error) {
				require.NoError(t, err)
			},
		},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			runAccessTokensWhiteListTestCase(t, &tc)
		})
	}
}

func TestAuthRepository_DeleteAccessTokenFromWhiteList(t *testing.T) {
	const token = "testAccessToken"

	tcs := map[string]authRepositoryTestCase{
		"SUCCESS": {
			reg: func(_ *gomock.Controller, conn *_cacheMock.MockConn, _ *_uuidMock.MockGenerator) {
				conn.EXPECT().Delete(gomock.Any(), formatAccessTokenIntoCacheKey(token)).
					Return(int64(1), nil)
			},
			cmd: func(repo *AuthRepository) error {
				return repo.DeleteAccessTokenFromWhitelist(context.Background(), token)
			},
			exp: func(err error) {
				require.NoError(t, err)
			},
		},
		"FAILURE": {
			reg: func(_ *gomock.Controller, conn *_cacheMock.MockConn, _ *_uuidMock.MockGenerator) {
				conn.EXPECT().Delete(gomock.Any(), formatAccessTokenIntoCacheKey(token)).
					Return(int64(0), errors.New(""))
			},
			cmd: func(repo *AuthRepository) error {
				return repo.DeleteAccessTokenFromWhitelist(context.Background(), token)
			},
			exp: func(err error) {
				require.Error(t, err)
			},
		},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			runAccessTokensWhiteListTestCase(t, &tc)
		})
	}
}

func runAccessTokensWhiteListTestCase(t *testing.T, tc *authRepositoryTestCase) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	whiteListConn := _cacheMock.NewMockConn(ctrl)
	idGen := _uuidMock.NewMockGenerator(ctrl)
	if tc.reg != nil {
		tc.reg(ctrl, whiteListConn, idGen)
	}

	logger := slog.NewLogger(stdslog.New(slogzap.Option{}.NewZapHandler()))
	repo := NewAuthRepository(nil, whiteListConn, nil, logger, idGen)

	var err error
	if tc.cmd != nil {
		err = tc.cmd(repo)
	}

	if tc.exp != nil {
		tc.exp(err)
	}
}

func TestAuthRepository_SaveRefreshTokenToBlackList(t *testing.T) {
	const token = "testRefreshToken"
	const expiration = time.Duration(-1)

	tcs := map[string]authRepositoryTestCase{
		"SUCCESS": {
			reg: func(_ *gomock.Controller, conn *_cacheMock.MockConn, _ *_uuidMock.MockGenerator) {
				conn.EXPECT().Set(gomock.Any(), formatRefreshTokenIntoCacheKey(token), token, expiration).
					Return(nil)
			},
			cmd: func(repo *AuthRepository) error {
				return repo.SaveRefreshTokenToBlacklist(context.Background(), token, expiration)
			},
			exp: func(err error) {
				require.NoError(t, err)
			},
		},
		"FAILURE": {
			reg: func(_ *gomock.Controller, conn *_cacheMock.MockConn, _ *_uuidMock.MockGenerator) {
				conn.EXPECT().Set(gomock.Any(), formatRefreshTokenIntoCacheKey(token), token, expiration).
					Return(errors.New(""))
			},
			cmd: func(repo *AuthRepository) error {
				return repo.SaveRefreshTokenToBlacklist(context.Background(), token, expiration)
			},
			exp: func(err error) {
				require.Error(t, err)
			},
		},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			runRefreshTokensBlackListTestCase(t, &tc)
		})
	}
}

func TestAuthRepository_FindRefreshTokenInBlackList(t *testing.T) {
	const token = "testRefreshToken"

	tcs := map[string]authRepositoryTestCase{
		"SUCCESS": {
			reg: func(_ *gomock.Controller, conn *_cacheMock.MockConn, _ *_uuidMock.MockGenerator) {
				conn.EXPECT().Exists(gomock.Any(), formatRefreshTokenIntoCacheKey(token)).
					Return(int64(1), nil)
			},
			cmd: func(repo *AuthRepository) error {
				found, err := repo.FindRefreshTokenInBlacklist(context.Background(), token)
				require.True(t, found)
				return err
			},
			exp: func(err error) {
				require.NoError(t, err)
			},
		},
		"FAILURE": {
			reg: func(_ *gomock.Controller, conn *_cacheMock.MockConn, _ *_uuidMock.MockGenerator) {
				conn.EXPECT().Exists(gomock.Any(), formatRefreshTokenIntoCacheKey(token)).
					Return(int64(0), errors.New(""))
			},
			cmd: func(repo *AuthRepository) error {
				found, err := repo.FindRefreshTokenInBlacklist(context.Background(), token)
				require.False(t, found)
				return err
			},
			exp: func(err error) {
				require.Error(t, err)
			},
		},
		"FAILURE not found": {
			reg: func(_ *gomock.Controller, conn *_cacheMock.MockConn, _ *_uuidMock.MockGenerator) {
				conn.EXPECT().Exists(gomock.Any(), formatRefreshTokenIntoCacheKey(token)).
					Return(int64(0), nil)
			},
			cmd: func(repo *AuthRepository) error {
				found, err := repo.FindRefreshTokenInBlacklist(context.Background(), token)
				require.False(t, found)
				return err
			},
			exp: func(err error) {
				require.NoError(t, err)
			},
		},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			runRefreshTokensBlackListTestCase(t, &tc)
		})
	}
}

func TestAuthRepository_DeleteRefreshTokenFromBlackList(t *testing.T) {
	const token = "testRefreshToken"

	tcs := map[string]authRepositoryTestCase{
		"SUCCESS": {
			reg: func(_ *gomock.Controller, conn *_cacheMock.MockConn, _ *_uuidMock.MockGenerator) {
				conn.EXPECT().Delete(gomock.Any(), formatRefreshTokenIntoCacheKey(token)).
					Return(int64(1), nil)
			},
			cmd: func(repo *AuthRepository) error {
				return repo.DeleteRefreshTokenFromBlacklist(context.Background(), token)
			},
			exp: func(err error) {
				require.NoError(t, err)
			},
		},
		"FAILURE": {
			reg: func(_ *gomock.Controller, conn *_cacheMock.MockConn, _ *_uuidMock.MockGenerator) {
				conn.EXPECT().Delete(gomock.Any(), formatRefreshTokenIntoCacheKey(token)).
					Return(int64(0), errors.New(""))
			},
			cmd: func(repo *AuthRepository) error {
				return repo.DeleteRefreshTokenFromBlacklist(context.Background(), token)
			},
			exp: func(err error) {
				require.Error(t, err)
			},
		},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			runRefreshTokensBlackListTestCase(t, &tc)
		})
	}
}

func runRefreshTokensBlackListTestCase(t *testing.T, tc *authRepositoryTestCase) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	blackListConn := _cacheMock.NewMockConn(ctrl)
	idGen := _uuidMock.NewMockGenerator(ctrl)
	if tc.reg != nil {
		tc.reg(ctrl, blackListConn, idGen)
	}

	logger := slog.NewLogger(stdslog.New(slogzap.Option{}.NewZapHandler()))
	repo := NewAuthRepository(nil, nil, blackListConn, logger, idGen)

	var err error
	if tc.cmd != nil {
		err = tc.cmd(repo)
	}

	if tc.exp != nil {
		tc.exp(err)
	}
}
