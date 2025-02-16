package redis

import (
	"context"
	"errors"
	_cacheMock "github.com/adanyl0v/pocket-ideas/mocks/pkg/cache"
	_uuidMock "github.com/adanyl0v/pocket-ideas/mocks/pkg/uuid"
	"github.com/adanyl0v/pocket-ideas/pkg/log/slog"
	"github.com/adanyl0v/pocket-ideas/pkg/proxerr"
	"github.com/golang/mock/gomock"
	slogzap "github.com/samber/slog-zap/v2"
	"github.com/stretchr/testify/require"
	stdslog "log/slog"
	"testing"
	"time"
)

type (
	accessTokensWhiteListTestCase struct {
		reg accessTokensWhiteListTestCaseRegister
		cmd authRepositoryTestCaseCommand
		exp authRepositoryTestCaseExpect
	}

	accessTokensWhiteListTestCaseRegister func(_ *gomock.Controller, conn *_cacheMock.MockConn, _ *_uuidMock.MockGenerator)
	authRepositoryTestCaseCommand         func(repo *AuthRepository) error
	authRepositoryTestCaseExpect          func(err error)
)

func TestAuthRepository_SaveAccessTokenToWhiteList(t *testing.T) {
	const token = "testAccessToken"
	const expiration = time.Duration(-1)

	tcs := map[string]accessTokensWhiteListTestCase{
		"SUCCESS": {
			reg: func(_ *gomock.Controller, conn *_cacheMock.MockConn, _ *_uuidMock.MockGenerator) {
				conn.EXPECT().Set(gomock.Any(), convertAccessTokenToWhiteListKey(token), token, expiration).
					Times(1).Return(nil)
			},
			cmd: func(repo *AuthRepository) error {
				return repo.SaveAccessTokenToWhiteList(context.Background(), token, expiration)
			},
			exp: func(err error) {
				require.NoError(t, err)
			},
		},
		"FAILURE": {
			reg: func(_ *gomock.Controller, conn *_cacheMock.MockConn, _ *_uuidMock.MockGenerator) {
				conn.EXPECT().Set(gomock.Any(), convertAccessTokenToWhiteListKey(token), token, expiration).
					Times(1).Return(errors.New(""))
			},
			cmd: func(repo *AuthRepository) error {
				return repo.SaveAccessTokenToWhiteList(context.Background(), token, expiration)
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

	tcs := map[string]accessTokensWhiteListTestCase{
		"SUCCESS": {
			reg: func(_ *gomock.Controller, conn *_cacheMock.MockConn, _ *_uuidMock.MockGenerator) {
				conn.EXPECT().Exists(gomock.Any(), convertAccessTokenToWhiteListKey(token)).
					Times(1).Return(int64(1), nil)
			},
			cmd: func(repo *AuthRepository) error {
				found, err := repo.FindAccessTokenInWhiteList(context.Background(), token)
				require.True(t, found)
				return err
			},
			exp: func(err error) {
				require.NoError(t, err)
			},
		},
		"FAILURE": {
			reg: func(_ *gomock.Controller, conn *_cacheMock.MockConn, _ *_uuidMock.MockGenerator) {
				conn.EXPECT().Exists(gomock.Any(), convertAccessTokenToWhiteListKey(token)).
					Times(1).Return(int64(0), errors.New(""))
			},
			cmd: func(repo *AuthRepository) error {
				found, err := repo.FindAccessTokenInWhiteList(context.Background(), token)
				require.False(t, found)
				return err
			},
			exp: func(err error) {
				require.Error(t, err)
			},
		},
		"FAILURE not found": {
			reg: func(_ *gomock.Controller, conn *_cacheMock.MockConn, _ *_uuidMock.MockGenerator) {
				conn.EXPECT().Exists(gomock.Any(), convertAccessTokenToWhiteListKey(token)).
					Times(1).Return(int64(0), nil)
			},
			cmd: func(repo *AuthRepository) error {
				found, err := repo.FindAccessTokenInWhiteList(context.Background(), token)
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

	tcs := map[string]accessTokensWhiteListTestCase{
		"SUCCESS": {
			reg: func(_ *gomock.Controller, conn *_cacheMock.MockConn, _ *_uuidMock.MockGenerator) {
				conn.EXPECT().Delete(gomock.Any(), convertAccessTokenToWhiteListKey(token)).
					Times(1).Return(nil)
			},
			cmd: func(repo *AuthRepository) error {
				return repo.DeleteAccessTokenFromWhiteList(context.Background(), token)
			},
			exp: func(err error) {
				require.NoError(t, err)
			},
		},
		"FAILURE": {
			reg: func(_ *gomock.Controller, conn *_cacheMock.MockConn, _ *_uuidMock.MockGenerator) {
				conn.EXPECT().Delete(gomock.Any(), convertAccessTokenToWhiteListKey(token)).
					Times(1).Return(errors.New(""))
			},
			cmd: func(repo *AuthRepository) error {
				return repo.DeleteAccessTokenFromWhiteList(context.Background(), token)
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

func runAccessTokensWhiteListTestCase(t *testing.T, tc *accessTokensWhiteListTestCase) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	conn := _cacheMock.NewMockConn(ctrl)
	idGen := _uuidMock.NewMockGenerator(ctrl)
	if tc.reg != nil {
		tc.reg(ctrl, conn, idGen)
	}

	logger := slog.NewLogger(stdslog.New(slogzap.Option{}.NewZapHandler()))
	repo := NewAuthRepository(nil, conn, nil, logger, idGen)

	var err error
	if tc.cmd != nil {
		err = tc.cmd(repo)
	}

	if err != nil {
		var pxErr proxerr.Error
		if errors.As(err, &pxErr) {
			err = pxErr.Unwrap()
		}
	}

	if tc.exp != nil {
		tc.exp(err)
	}
}
