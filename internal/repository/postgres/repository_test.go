package postgres

import (
	"context"
	"errors"
	_dbMock "github.com/adanyl0v/pocket-ideas/mocks/pkg/database"
	"github.com/adanyl0v/pocket-ideas/pkg/log/slog"
	"github.com/golang/mock/gomock"
	slogzap "github.com/samber/slog-zap/v2"
	"github.com/stretchr/testify/require"
	stdslog "log/slog"
	"testing"
)

func TestRepository_Begin(t *testing.T) {
	tcs := map[string]struct {
		reg func(_ *gomock.Controller, conn *_dbMock.MockConn)
		cmd func(repo *Repository) error
		exp func(err error)
	}{
		"SUCCESS": {
			reg: func(_ *gomock.Controller, conn *_dbMock.MockConn) {
				conn.EXPECT().Begin(gomock.Any()).Return(nil, nil)
			},
			cmd: func(repo *Repository) error {
				_, err := repo.Begin(context.Background())
				return err
			},
			exp: func(err error) {
				require.NoError(t, err)
			},
		},
		"FAILED": {
			reg: func(_ *gomock.Controller, conn *_dbMock.MockConn) {
				conn.EXPECT().Begin(gomock.Any()).Return(nil, errors.New(""))
			},
			cmd: func(repo *Repository) error {
				_, err := repo.Begin(context.Background())
				return err
			},
			exp: func(err error) {
				require.Error(t, err)
			},
		},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			conn := _dbMock.NewMockConn(ctrl)
			if tc.reg != nil {
				tc.reg(ctrl, conn)
			}

			logger := slog.NewLogger(stdslog.New(slogzap.Option{}.NewZapHandler()))
			repo := &Repository{
				conn:   conn,
				logger: logger,
			}

			var err error
			if tc.cmd != nil {
				err = tc.cmd(repo)
			}

			if tc.exp != nil {
				tc.exp(err)
			}
		})
	}
}
