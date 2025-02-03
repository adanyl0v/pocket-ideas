package postgres

import (
	"context"
	"errors"
	"github.com/adanyl0v/pocket-ideas/internal/domain"
	"github.com/adanyl0v/pocket-ideas/internal/repository"
	_pgdbMock "github.com/adanyl0v/pocket-ideas/mocks/pkg/database/postgres"
	_uuidMock "github.com/adanyl0v/pocket-ideas/mocks/pkg/uuid"
	pgdb "github.com/adanyl0v/pocket-ideas/pkg/database/postgres"
	"github.com/adanyl0v/pocket-ideas/pkg/log/slog"
	"github.com/adanyl0v/pocket-ideas/pkg/proxerr"
	"github.com/golang/mock/gomock"
	slogzap "github.com/samber/slog-zap/v2"
	"github.com/stretchr/testify/require"
	stdslog "log/slog"
	"testing"
)

type (
	userTc struct {
		register userTcRegister
		command  userTcCommand
		expect   userTcExpect
	}

	userTcRegister func(ctrl *gomock.Controller, conn *_pgdbMock.MockConn, idGen *_uuidMock.MockGenerator)
	userTcCommand  func(repo *UserRepository) error
	userTcExpect   func(err error)
)

func TestNewUserRepository(t *testing.T) {
	NewUserRepository(nil, nil, nil)
}

func TestUserRepository_Begin(t *testing.T) {
	testCases := map[string]userTc{
		"success": {
			register: func(ctrl *gomock.Controller, conn *_pgdbMock.MockConn, _ *_uuidMock.MockGenerator) {
				tx := _pgdbMock.NewMockTx(ctrl)
				conn.EXPECT().Begin(gomock.Any()).Times(1).Return(tx, nil)
			},
			command: func(repo *UserRepository) error {
				_, err := repo.Begin(context.Background())
				return err
			},
			expect: func(err error) {
				require.NoError(t, err)
			},
		},
		"failure": {
			register: func(_ *gomock.Controller, conn *_pgdbMock.MockConn, _ *_uuidMock.MockGenerator) {
				conn.EXPECT().Begin(gomock.Any()).Times(1).Return(nil, errors.New(""))
			},
			command: func(repo *UserRepository) error {
				_, err := repo.Begin(context.Background())
				return err
			},
			expect: func(err error) {
				require.NotNil(t, err)
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			runUserTestCase(t, &tc)
		})
	}
}

func TestUserRepository_WithTx(t *testing.T) {
	tcs := map[string]userTc{
		"success": {
			command: func(repo *UserRepository) error {
				r := repo.WithTx(new(_pgdbMock.MockTx))
				require.NotNil(t, r)
				return nil
			},
		},
		"failure": {
			command: func(repo *UserRepository) error {
				r := repo.WithTx(repository.Tx(nil))
				require.Nil(t, r)
				return nil
			},
		},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			runUserTestCase(t, &tc)
		})
	}
}

func TestUserRepository_Create(t *testing.T) {
	testCases := map[string]userTc{
		"success": {
			register: func(_ *gomock.Controller, conn *_pgdbMock.MockConn, idGen *_uuidMock.MockGenerator) {
				idGen.EXPECT().NewV7().Times(1).Return("", nil)
				conn.EXPECT().Exec(gomock.Any(), insertUserQuery, gomock.Any(), gomock.Any(), gomock.Any(),
					gomock.Any(), gomock.Any(), gomock.Any()).Times(1).Return(nil)
			},
			command: func(repo *UserRepository) error {
				return repo.Create(context.Background(), new(domain.User))
			},
			expect: func(err error) {
				require.Nil(t, err)
			},
		},
		"failed user id generation": {
			register: func(_ *gomock.Controller, _ *_pgdbMock.MockConn, idGen *_uuidMock.MockGenerator) {
				idGen.EXPECT().NewV7().Times(1).Return("", errGenerateUserId)
			},
			command: func(repo *UserRepository) error {
				return repo.Create(context.Background(), new(domain.User))
			},
			expect: func(err error) {
				require.Equal(t, errGenerateUserId, err)
			},
		},
		"user already exists": {
			register: func(_ *gomock.Controller, conn *_pgdbMock.MockConn, idGen *_uuidMock.MockGenerator) {
				idGen.EXPECT().NewV7().Times(1).Return("", nil)
				conn.EXPECT().Exec(gomock.Any(), insertUserQuery, gomock.Any(), gomock.Any(), gomock.Any(),
					gomock.Any(), gomock.Any(), gomock.Any()).Times(1).
					Return(proxerr.New(pgdb.ErrUniqueViolation, ""))
			},
			command: func(repo *UserRepository) error {
				return repo.Create(context.Background(), new(domain.User))
			},
			expect: func(err error) {
				require.Equal(t, ErrUserExists, err)
			},
		},
		"failure": {
			register: func(ctrl *gomock.Controller, conn *_pgdbMock.MockConn, idGen *_uuidMock.MockGenerator) {
				idGen.EXPECT().NewV7().Times(1).Return("", nil)
				conn.EXPECT().Exec(gomock.Any(), insertUserQuery, gomock.Any(), gomock.Any(), gomock.Any(),
					gomock.Any(), gomock.Any(), gomock.Any()).Times(1).
					Return(errors.New(""))
			},
			command: func(repo *UserRepository) error {
				return repo.Create(context.Background(), new(domain.User))
			},
			expect: func(err error) {
				require.NotNil(t, err)
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			runUserTestCase(t, &tc)
		})
	}
}

func TestUserRepository_GetByID(t *testing.T) {
	testCases := map[string]userTc{
		"success": {
			register: func(ctrl *gomock.Controller, conn *_pgdbMock.MockConn, _ *_uuidMock.MockGenerator) {
				row := _pgdbMock.NewMockRow(ctrl)
				row.EXPECT().Scan(gomock.Any(), gomock.Any(), gomock.Any(),
					gomock.Any(), gomock.Any()).Times(1).Return(nil)
				conn.EXPECT().QueryRow(gomock.Any(), getUserByIdQuery, gomock.Any()).Times(1).Return(row)
			},
			command: func(repo *UserRepository) error {
				_, err := repo.GetByID(context.Background(), "")
				return err
			},
			expect: func(err error) {
				require.Nil(t, err)
			},
		},
		"user not found": {
			register: func(ctrl *gomock.Controller, conn *_pgdbMock.MockConn, idGen *_uuidMock.MockGenerator) {
				row := _pgdbMock.NewMockRow(ctrl)
				row.EXPECT().Scan(gomock.Any(), gomock.Any(), gomock.Any(),
					gomock.Any(), gomock.Any()).Times(1).Return(proxerr.New(pgdb.ErrNoRows, ""))
				conn.EXPECT().QueryRow(gomock.Any(), getUserByIdQuery, gomock.Any()).Times(1).Return(row)
			},
			command: func(repo *UserRepository) error {
				_, err := repo.GetByID(context.Background(), "")
				return err
			},
			expect: func(err error) {
				require.Equal(t, ErrUserNotFound, err)
			},
		},
		"failure": {
			register: func(ctrl *gomock.Controller, conn *_pgdbMock.MockConn, _ *_uuidMock.MockGenerator) {
				row := _pgdbMock.NewMockRow(ctrl)
				row.EXPECT().Scan(gomock.Any(), gomock.Any(), gomock.Any(),
					gomock.Any(), gomock.Any()).Times(1).Return(errors.New(""))
				conn.EXPECT().QueryRow(gomock.Any(), getUserByIdQuery, gomock.Any()).Times(1).Return(row)
			},
			command: func(repo *UserRepository) error {
				_, err := repo.GetByID(context.Background(), "")
				return err
			},
			expect: func(err error) {
				require.NotNil(t, err)
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			runUserTestCase(t, &tc)
		})
	}
}

func TestUserRepository_GetByEmail(t *testing.T) {
	testCases := map[string]userTc{
		"success": {
			register: func(ctrl *gomock.Controller, conn *_pgdbMock.MockConn, _ *_uuidMock.MockGenerator) {
				row := _pgdbMock.NewMockRow(ctrl)
				row.EXPECT().Scan(gomock.Any(), gomock.Any(), gomock.Any(),
					gomock.Any(), gomock.Any()).Times(1).Return(nil)
				conn.EXPECT().QueryRow(gomock.Any(), getUserByEmailQuery, gomock.Any()).Times(1).Return(row)
			},
			command: func(repo *UserRepository) error {
				_, err := repo.GetByEmail(context.Background(), "")
				return err
			},
			expect: func(err error) {
				require.Nil(t, err)
			},
		},
		"user not found": {
			register: func(ctrl *gomock.Controller, conn *_pgdbMock.MockConn, idGen *_uuidMock.MockGenerator) {
				row := _pgdbMock.NewMockRow(ctrl)
				row.EXPECT().Scan(gomock.Any(), gomock.Any(), gomock.Any(),
					gomock.Any(), gomock.Any()).Times(1).Return(proxerr.New(pgdb.ErrNoRows, ""))
				conn.EXPECT().QueryRow(gomock.Any(), getUserByEmailQuery, gomock.Any()).Times(1).Return(row)
			},
			command: func(repo *UserRepository) error {
				_, err := repo.GetByEmail(context.Background(), "")
				return err
			},
			expect: func(err error) {
				require.Equal(t, ErrUserNotFound, err)
			},
		},
		"failure": {
			register: func(ctrl *gomock.Controller, conn *_pgdbMock.MockConn, _ *_uuidMock.MockGenerator) {
				row := _pgdbMock.NewMockRow(ctrl)
				row.EXPECT().Scan(gomock.Any(), gomock.Any(), gomock.Any(),
					gomock.Any(), gomock.Any()).Times(1).Return(errors.New(""))
				conn.EXPECT().QueryRow(gomock.Any(), getUserByEmailQuery, gomock.Any()).Times(1).Return(row)
			},
			command: func(repo *UserRepository) error {
				_, err := repo.GetByEmail(context.Background(), "")
				return err
			},
			expect: func(err error) {
				require.NotNil(t, err)
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			runUserTestCase(t, &tc)
		})
	}
}

func TestUserRepository_SelectAll(t *testing.T) {
	testCases := map[string]userTc{
		"success": {
			register: func(ctrl *gomock.Controller, conn *_pgdbMock.MockConn, _ *_uuidMock.MockGenerator) {
				rows := _pgdbMock.NewMockRows(ctrl)
				rows.EXPECT().Close().Times(1)

				i, n := -1, 5
				rows.EXPECT().Next().Times(n + 1).DoAndReturn(func() bool {
					i++
					return i < n
				})

				rows.EXPECT().Scan(gomock.Any(), gomock.Any(), gomock.Any(),
					gomock.Any(), gomock.Any(), gomock.Any()).Times(n).Return(nil)
				conn.EXPECT().Query(gomock.Any(), selectAllUsersQuery).Times(1).Return(rows, nil)
			},
			command: func(repo *UserRepository) error {
				_, err := repo.SelectAll(context.Background())
				return err
			},
			expect: func(err error) {
				require.Nil(t, err)
			},
		},
		"failed selection": {
			register: func(_ *gomock.Controller, conn *_pgdbMock.MockConn, _ *_uuidMock.MockGenerator) {
				conn.EXPECT().Query(gomock.Any(), selectAllUsersQuery).Times(1).Return(nil, errors.New(""))
			},
			command: func(repo *UserRepository) error {
				_, err := repo.SelectAll(context.Background())
				return err
			},
			expect: func(err error) {
				require.NotNil(t, err)
			},
		},
		"no users found": {
			register: func(ctrl *gomock.Controller, conn *_pgdbMock.MockConn, _ *_uuidMock.MockGenerator) {
				rows := _pgdbMock.NewMockRows(ctrl)
				rows.EXPECT().Close().Times(1)
				rows.EXPECT().Next().Times(1).Return(false)
				rows.EXPECT().Err().Times(1).Return(nil)
				conn.EXPECT().Query(gomock.Any(), selectAllUsersQuery).Times(1).Return(rows, nil)
			},
			command: func(repo *UserRepository) error {
				_, err := repo.SelectAll(context.Background())
				return err
			},
			expect: func(err error) {
				require.Equal(t, ErrNoUsersFound, err)
			},
		},
		"no users selected": {
			register: func(ctrl *gomock.Controller, conn *_pgdbMock.MockConn, _ *_uuidMock.MockGenerator) {
				rows := _pgdbMock.NewMockRows(ctrl)
				rows.EXPECT().Close().Times(1)
				rows.EXPECT().Next().Times(1).Return(false)
				rows.EXPECT().Err().Times(1).Return(errors.New(""))
				conn.EXPECT().Query(gomock.Any(), selectAllUsersQuery).Times(1).Return(rows, nil)
			},
			command: func(repo *UserRepository) error {
				_, err := repo.SelectAll(context.Background())
				return err
			},
			expect: func(err error) {
				require.NotNil(t, err)
			},
		},
		"failed scan": {
			register: func(ctrl *gomock.Controller, conn *_pgdbMock.MockConn, _ *_uuidMock.MockGenerator) {
				rows := _pgdbMock.NewMockRows(ctrl)
				rows.EXPECT().Close().Times(1)
				rows.EXPECT().Next().Times(1).Return(true)
				rows.EXPECT().Scan(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
					gomock.Any(), gomock.Any()).Times(1).Return(errors.New(""))
				conn.EXPECT().Query(gomock.Any(), selectAllUsersQuery).Times(1).Return(rows, nil)
			},
			command: func(repo *UserRepository) error {
				_, err := repo.SelectAll(context.Background())
				return err
			},
			expect: func(err error) {
				require.NotNil(t, err)
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			runUserTestCase(t, &tc)
		})
	}
}

func TestUserRepository_SelectByName(t *testing.T) {
	testCases := map[string]userTc{
		"success": {
			register: func(ctrl *gomock.Controller, conn *_pgdbMock.MockConn, _ *_uuidMock.MockGenerator) {
				rows := _pgdbMock.NewMockRows(ctrl)
				rows.EXPECT().Close().Times(1)

				i, n := -1, 5
				rows.EXPECT().Next().Times(n + 1).DoAndReturn(func() bool {
					i++
					return i < n
				})

				rows.EXPECT().Scan(gomock.Any(), gomock.Any(), gomock.Any(),
					gomock.Any(), gomock.Any()).Times(n).Return(nil)
				conn.EXPECT().Query(gomock.Any(), selectUsersByNameQuery, "").Times(1).Return(rows, nil)
			},
			command: func(repo *UserRepository) error {
				_, err := repo.SelectByName(context.Background(), "")
				return err
			},
			expect: func(err error) {
				require.Nil(t, err)
			},
		},
		"failed selection": {
			register: func(_ *gomock.Controller, conn *_pgdbMock.MockConn, _ *_uuidMock.MockGenerator) {
				conn.EXPECT().Query(gomock.Any(), selectUsersByNameQuery, "").
					Times(1).Return(nil, errors.New(""))
			},
			command: func(repo *UserRepository) error {
				_, err := repo.SelectByName(context.Background(), "")
				return err
			},
			expect: func(err error) {
				require.NotNil(t, err)
			},
		},
		"no users found": {
			register: func(ctrl *gomock.Controller, conn *_pgdbMock.MockConn, _ *_uuidMock.MockGenerator) {
				rows := _pgdbMock.NewMockRows(ctrl)
				rows.EXPECT().Close().Times(1)
				rows.EXPECT().Next().Times(1).Return(false)
				rows.EXPECT().Err().Times(1).Return(nil)
				conn.EXPECT().Query(gomock.Any(), selectUsersByNameQuery, "").Times(1).Return(rows, nil)
			},
			command: func(repo *UserRepository) error {
				_, err := repo.SelectByName(context.Background(), "")
				return err
			},
			expect: func(err error) {
				require.Equal(t, ErrNoUsersFound, err)
			},
		},
		"no users selected": {
			register: func(ctrl *gomock.Controller, conn *_pgdbMock.MockConn, _ *_uuidMock.MockGenerator) {
				rows := _pgdbMock.NewMockRows(ctrl)
				rows.EXPECT().Close().Times(1)
				rows.EXPECT().Next().Times(1).Return(false)
				rows.EXPECT().Err().Times(1).Return(errors.New(""))
				conn.EXPECT().Query(gomock.Any(), selectUsersByNameQuery, "").Times(1).Return(rows, nil)
			},
			command: func(repo *UserRepository) error {
				_, err := repo.SelectByName(context.Background(), "")
				return err
			},
			expect: func(err error) {
				require.NotNil(t, err)
			},
		},
		"failed scan": {
			register: func(ctrl *gomock.Controller, conn *_pgdbMock.MockConn, _ *_uuidMock.MockGenerator) {
				rows := _pgdbMock.NewMockRows(ctrl)
				rows.EXPECT().Close().Times(1)
				rows.EXPECT().Next().Times(1).Return(true)
				rows.EXPECT().Scan(gomock.Any(), gomock.Any(), gomock.Any(),
					gomock.Any(), gomock.Any()).Times(1).Return(errors.New(""))
				conn.EXPECT().Query(gomock.Any(), selectUsersByNameQuery, "").Times(1).Return(rows, nil)
			},
			command: func(repo *UserRepository) error {
				_, err := repo.SelectByName(context.Background(), "")
				return err
			},
			expect: func(err error) {
				require.NotNil(t, err)
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			runUserTestCase(t, &tc)
		})
	}
}

// runUserTestCase should be called by [t.Run] to enable
// the Run gutter icon (at least in Goland from JetBrains)
//
// [t.Run]: https://pkg.go.dev/testing#T.Run
func runUserTestCase(t *testing.T, tc *userTc) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	conn := _pgdbMock.NewMockConn(ctrl)
	idGen := _uuidMock.NewMockGenerator(ctrl)
	if tc.register != nil {
		tc.register(ctrl, conn, idGen)
	}

	logger := slog.NewLogger(stdslog.New(slogzap.Option{}.NewZapHandler()))
	repo := NewUserRepository(conn, logger, idGen)

	var err error
	if tc.command != nil {
		err = tc.command(repo)
	}

	if err != nil {
		var pxErr proxerr.Error
		if errors.As(err, &pxErr) {
			err = pxErr.Unwrap()
		}
	}

	if tc.expect != nil {
		tc.expect(err)
	}
}
