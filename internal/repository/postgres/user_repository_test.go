package postgres

import (
	"context"
	"errors"
	"github.com/adanyl0v/pocket-ideas/internal/domain"
	_dbMock "github.com/adanyl0v/pocket-ideas/mocks/pkg/database"
	_uuidMock "github.com/adanyl0v/pocket-ideas/mocks/pkg/uuid"
	"github.com/adanyl0v/pocket-ideas/pkg/database"
	"github.com/adanyl0v/pocket-ideas/pkg/log/slog"
	"github.com/adanyl0v/pocket-ideas/pkg/proxerr"
	"github.com/golang/mock/gomock"
	slogzap "github.com/samber/slog-zap/v2"
	"github.com/stretchr/testify/require"
	stdslog "log/slog"
	"testing"
)

type (
	userTestCaseRegister func(_ *gomock.Controller, conn *_dbMock.MockConn, _ *_uuidMock.MockGenerator)
	userTestCaseCommand  func(repo *UserRepository) error
	userTestCaseExpect   func(err error)

	userTestCase struct {
		reg userTestCaseRegister
		cmd userTestCaseCommand
		exp userTestCaseExpect
	}
)

func TestUserRepository_WithTx(t *testing.T) {
	tcs := map[string]userTestCase{
		"SUCCESS": {
			cmd: func(repo *UserRepository) error {
				r := repo.WithTx(new(_dbMock.MockTx))
				require.NotNil(t, r)
				return nil
			},
		},
		"FAILED": {
			cmd: func(repo *UserRepository) error {
				r := repo.WithTx(nil)
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

func TestUserRepository_Save(t *testing.T) {
	tcs := map[string]userTestCase{
		"SUCCESS": {
			reg: func(_ *gomock.Controller, conn *_dbMock.MockConn, idGen *_uuidMock.MockGenerator) {
				idGen.EXPECT().NewV7().Times(1).Return("", nil)
				conn.EXPECT().Execute(gomock.Any(), qInsertUser, gomock.Any(), gomock.Any(), gomock.Any(),
					gomock.Any(), gomock.Any(), gomock.Any()).Times(1).Return(nil, nil)
			},
			cmd: func(repo *UserRepository) error {
				return repo.Save(context.Background(), new(domain.User))
			},
			exp: func(err error) {
				require.NoError(t, err)
			},
		},
		"FAILED id generation": {
			reg: func(_ *gomock.Controller, _ *_dbMock.MockConn, idGen *_uuidMock.MockGenerator) {
				idGen.EXPECT().NewV7().Times(1).Return("", errors.New(""))
			},
			cmd: func(repo *UserRepository) error {
				return repo.Save(context.Background(), new(domain.User))
			},
			exp: func(err error) {
				require.Error(t, err)
			},
		},
		"FAILED user already exists": {
			reg: func(_ *gomock.Controller, conn *_dbMock.MockConn, idGen *_uuidMock.MockGenerator) {
				idGen.EXPECT().NewV7().Times(1).Return("", nil)
				conn.EXPECT().Execute(gomock.Any(), qInsertUser, gomock.Any(), gomock.Any(), gomock.Any(),
					gomock.Any(), gomock.Any(), gomock.Any()).Times(1).Return(nil,
					proxerr.New(database.ErrUniqueViolation, ""))
			},
			cmd: func(repo *UserRepository) error {
				return repo.Save(context.Background(), new(domain.User))
			},
			exp: func(err error) {
				require.Equal(t, ErrUserAlreadyExists, err)
			},
		},
		"FAILED user field must not be empty": {
			reg: func(_ *gomock.Controller, conn *_dbMock.MockConn, idGen *_uuidMock.MockGenerator) {
				idGen.EXPECT().NewV7().Times(1).Return("", nil)
				conn.EXPECT().Execute(gomock.Any(), qInsertUser, gomock.Any(), gomock.Any(), gomock.Any(),
					gomock.Any(), gomock.Any(), gomock.Any()).Times(1).Return(nil,
					proxerr.New(database.ErrNotNullViolation, ""))
			},
			cmd: func(repo *UserRepository) error {
				return repo.Save(context.Background(), new(domain.User))
			},
			exp: func(err error) {
				require.Equal(t, ErrUserFieldMustNotBeEmpty, err)
			},
		},
		"FAILED": {
			reg: func(_ *gomock.Controller, conn *_dbMock.MockConn, idGen *_uuidMock.MockGenerator) {
				idGen.EXPECT().NewV7().Times(1).Return("", nil)
				conn.EXPECT().Execute(gomock.Any(), qInsertUser, gomock.Any(), gomock.Any(), gomock.Any(),
					gomock.Any(), gomock.Any(), gomock.Any()).Times(1).Return(nil, errors.New(""))
			},
			cmd: func(repo *UserRepository) error {
				return repo.Save(context.Background(), new(domain.User))
			},
			exp: func(err error) {
				require.Error(t, err)
			},
		},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			runUserTestCase(t, &tc)
		})
	}
}

func TestUserRepository_FindById(t *testing.T) {
	const id = "0194f574-5a05-7e68-91d6-d30f1d81869c"
	tcs := map[string]userTestCase{
		"SUCCESS": {
			reg: func(ctrl *gomock.Controller, conn *_dbMock.MockConn, _ *_uuidMock.MockGenerator) {
				row := _dbMock.NewMockRow(ctrl)
				row.EXPECT().Scan(gomock.Any(), gomock.Any(), gomock.Any(),
					gomock.Any(), gomock.Any()).Times(1).Return(nil)

				conn.EXPECT().QueryRow(gomock.Any(), qFindUserById, id).Times(1).Return(row)
			},
			cmd: func(repo *UserRepository) error {
				_, err := repo.FindById(context.Background(), id)
				return err
			},
			exp: func(err error) {
				require.NoError(t, err)
			},
		},
		"FAILED user not found": {
			reg: func(ctrl *gomock.Controller, conn *_dbMock.MockConn, _ *_uuidMock.MockGenerator) {
				row := _dbMock.NewMockRow(ctrl)
				row.EXPECT().Scan(gomock.Any(), gomock.Any(), gomock.Any(),
					gomock.Any(), gomock.Any()).Times(1).Return(proxerr.New(database.ErrNoRows, ""))

				conn.EXPECT().QueryRow(gomock.Any(), qFindUserById, id).Times(1).Return(row)
			},
			cmd: func(repo *UserRepository) error {
				_, err := repo.FindById(context.Background(), id)
				return err
			},
			exp: func(err error) {
				require.Equal(t, ErrUserNotFound, err)
			},
		},
		"FAILED": {
			reg: func(ctrl *gomock.Controller, conn *_dbMock.MockConn, _ *_uuidMock.MockGenerator) {
				row := _dbMock.NewMockRow(ctrl)
				row.EXPECT().Scan(gomock.Any(), gomock.Any(), gomock.Any(),
					gomock.Any(), gomock.Any()).Times(1).Return(errors.New(""))

				conn.EXPECT().QueryRow(gomock.Any(), qFindUserById, id).Times(1).Return(row)
			},
			cmd: func(repo *UserRepository) error {
				_, err := repo.FindById(context.Background(), id)
				return err
			},
			exp: func(err error) {
				require.Error(t, err)
			},
		},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			runUserTestCase(t, &tc)
		})
	}
}

func TestUserRepository_FindByEmail(t *testing.T) {
	const email = "user@example.com"
	tcs := map[string]userTestCase{
		"SUCCESS": {
			reg: func(ctrl *gomock.Controller, conn *_dbMock.MockConn, _ *_uuidMock.MockGenerator) {
				row := _dbMock.NewMockRow(ctrl)
				row.EXPECT().Scan(gomock.Any(), gomock.Any(), gomock.Any(),
					gomock.Any(), gomock.Any()).Times(1).Return(nil)

				conn.EXPECT().QueryRow(gomock.Any(), qFindUserByEmail, email).Times(1).Return(row)
			},
			cmd: func(repo *UserRepository) error {
				_, err := repo.FindByEmail(context.Background(), email)
				return err
			},
			exp: func(err error) {
				require.NoError(t, err)
			},
		},
		"FAILED user not found": {
			reg: func(ctrl *gomock.Controller, conn *_dbMock.MockConn, _ *_uuidMock.MockGenerator) {
				row := _dbMock.NewMockRow(ctrl)
				row.EXPECT().Scan(gomock.Any(), gomock.Any(), gomock.Any(),
					gomock.Any(), gomock.Any()).Times(1).Return(proxerr.New(database.ErrNoRows, ""))

				conn.EXPECT().QueryRow(gomock.Any(), qFindUserByEmail, email).Times(1).Return(row)
			},
			cmd: func(repo *UserRepository) error {
				_, err := repo.FindByEmail(context.Background(), email)
				return err
			},
			exp: func(err error) {
				require.Equal(t, ErrUserNotFound, err)
			},
		},
		"FAILED": {
			reg: func(ctrl *gomock.Controller, conn *_dbMock.MockConn, _ *_uuidMock.MockGenerator) {
				row := _dbMock.NewMockRow(ctrl)
				row.EXPECT().Scan(gomock.Any(), gomock.Any(), gomock.Any(),
					gomock.Any(), gomock.Any()).Times(1).Return(errors.New(""))

				conn.EXPECT().QueryRow(gomock.Any(), qFindUserByEmail, email).Times(1).Return(row)
			},
			cmd: func(repo *UserRepository) error {
				_, err := repo.FindByEmail(context.Background(), email)
				return err
			},
			exp: func(err error) {
				require.Error(t, err)
			},
		},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			runUserTestCase(t, &tc)
		})
	}
}

func TestUserRepository_FindAll(t *testing.T) {
	tcs := map[string]userTestCase{
		"SUCCESS": {
			reg: func(ctrl *gomock.Controller, conn *_dbMock.MockConn, _ *_uuidMock.MockGenerator) {
				rows := _dbMock.NewMockRows(ctrl)
				rows.EXPECT().Close().Times(1)
				rows.EXPECT().Err().Return(nil).Times(1)

				var n, i = 5, 0
				rows.EXPECT().Next().Times(n).DoAndReturn(func() bool {
					i++
					return i < n
				})
				rows.EXPECT().Scan(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Times(n - 1).Return(nil)

				conn.EXPECT().Query(gomock.Any(), qFindAllUsers).Times(1).Return(rows, nil)
			},
			cmd: func(repo *UserRepository) error {
				_, err := repo.FindAll(context.Background())
				return err
			},
			exp: func(err error) {
				require.NoError(t, err)
			},
		},
		"FAILED query execution": {
			reg: func(ctrl *gomock.Controller, conn *_dbMock.MockConn, _ *_uuidMock.MockGenerator) {
				conn.EXPECT().Query(gomock.Any(), qFindAllUsers).Times(1).Return(nil, errors.New(""))
			},
			cmd: func(repo *UserRepository) error {
				_, err := repo.FindAll(context.Background())
				return err
			},
			exp: func(err error) {
				require.Error(t, err)
			},
		},
		"FAILED model scan": {
			reg: func(ctrl *gomock.Controller, conn *_dbMock.MockConn, _ *_uuidMock.MockGenerator) {
				rows := _dbMock.NewMockRows(ctrl)
				rows.EXPECT().Close().Times(1)
				rows.EXPECT().Next().Times(1).Return(true)
				rows.EXPECT().Scan(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1).Return(errors.New(""))

				conn.EXPECT().Query(gomock.Any(), qFindAllUsers).Times(1).Return(rows, nil)
			},
			cmd: func(repo *UserRepository) error {
				_, err := repo.FindAll(context.Background())
				return err
			},
			exp: func(err error) {
				require.Error(t, err)
			},
		},
		"FAILED user not found": {
			reg: func(ctrl *gomock.Controller, conn *_dbMock.MockConn, _ *_uuidMock.MockGenerator) {
				rows := _dbMock.NewMockRows(ctrl)
				rows.EXPECT().Close().Times(1)
				rows.EXPECT().Next().Times(1).Return(false)
				rows.EXPECT().Err().Return(proxerr.New(database.ErrNoRows, "")).Times(1)

				conn.EXPECT().Query(gomock.Any(), qFindAllUsers).Times(1).Return(rows, nil)
			},
			cmd: func(repo *UserRepository) error {
				_, err := repo.FindAll(context.Background())
				return err
			},
			exp: func(err error) {
				require.Equal(t, ErrUserNotFound, err)
			},
		},
		"FAILED": {
			reg: func(ctrl *gomock.Controller, conn *_dbMock.MockConn, _ *_uuidMock.MockGenerator) {
				rows := _dbMock.NewMockRows(ctrl)
				rows.EXPECT().Close().Times(1)
				rows.EXPECT().Next().Times(1).Return(false)
				rows.EXPECT().Err().Return(errors.New("")).Times(1)

				conn.EXPECT().Query(gomock.Any(), qFindAllUsers).Times(1).Return(rows, nil)
			},
			cmd: func(repo *UserRepository) error {
				_, err := repo.FindAll(context.Background())
				return err
			},
			exp: func(err error) {
				require.Error(t, err)
			},
		},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			runUserTestCase(t, &tc)
		})
	}
}

func TestUserRepository_FindByName(t *testing.T) {
	const name = "user"
	tcs := map[string]userTestCase{
		"SUCCESS": {
			reg: func(ctrl *gomock.Controller, conn *_dbMock.MockConn, _ *_uuidMock.MockGenerator) {
				rows := _dbMock.NewMockRows(ctrl)
				rows.EXPECT().Close().Times(1)
				rows.EXPECT().Err().Return(nil).Times(1)

				var n, i = 5, 0
				rows.EXPECT().Next().Times(n).DoAndReturn(func() bool {
					i++
					return i < n
				})
				rows.EXPECT().Scan(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Times(n - 1).Return(nil)

				conn.EXPECT().Query(gomock.Any(), qFindUsersByName, name).Times(1).Return(rows, nil)
			},
			cmd: func(repo *UserRepository) error {
				_, err := repo.FindByName(context.Background(), name)
				return err
			},
			exp: func(err error) {
				require.NoError(t, err)
			},
		},
		"FAILED query execution": {
			reg: func(ctrl *gomock.Controller, conn *_dbMock.MockConn, _ *_uuidMock.MockGenerator) {
				conn.EXPECT().Query(gomock.Any(), qFindUsersByName, name).Times(1).Return(nil, errors.New(""))
			},
			cmd: func(repo *UserRepository) error {
				_, err := repo.FindByName(context.Background(), name)
				return err
			},
			exp: func(err error) {
				require.Error(t, err)
			},
		},
		"FAILED model scan": {
			reg: func(ctrl *gomock.Controller, conn *_dbMock.MockConn, _ *_uuidMock.MockGenerator) {
				rows := _dbMock.NewMockRows(ctrl)
				rows.EXPECT().Close().Times(1)
				rows.EXPECT().Next().Times(1).Return(true)
				rows.EXPECT().Scan(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1).Return(errors.New(""))

				conn.EXPECT().Query(gomock.Any(), qFindUsersByName, name).Times(1).Return(rows, nil)
			},
			cmd: func(repo *UserRepository) error {
				_, err := repo.FindByName(context.Background(), name)
				return err
			},
			exp: func(err error) {
				require.Error(t, err)
			},
		},
		"FAILED user not found": {
			reg: func(ctrl *gomock.Controller, conn *_dbMock.MockConn, _ *_uuidMock.MockGenerator) {
				rows := _dbMock.NewMockRows(ctrl)
				rows.EXPECT().Close().Times(1)
				rows.EXPECT().Next().Times(1).Return(false)
				rows.EXPECT().Err().Return(proxerr.New(database.ErrNoRows, "")).Times(1)

				conn.EXPECT().Query(gomock.Any(), qFindUsersByName, name).Times(1).Return(rows, nil)
			},
			cmd: func(repo *UserRepository) error {
				_, err := repo.FindByName(context.Background(), name)
				return err
			},
			exp: func(err error) {
				require.Equal(t, ErrUserNotFound, err)
			},
		},
		"FAILED": {
			reg: func(ctrl *gomock.Controller, conn *_dbMock.MockConn, _ *_uuidMock.MockGenerator) {
				rows := _dbMock.NewMockRows(ctrl)
				rows.EXPECT().Close().Times(1)
				rows.EXPECT().Next().Times(1).Return(false)
				rows.EXPECT().Err().Return(errors.New("")).Times(1)

				conn.EXPECT().Query(gomock.Any(), qFindUsersByName, name).Times(1).Return(rows, nil)
			},
			cmd: func(repo *UserRepository) error {
				_, err := repo.FindByName(context.Background(), name)
				return err
			},
			exp: func(err error) {
				require.Error(t, err)
			},
		},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			runUserTestCase(t, &tc)
		})
	}
}

func TestUserRepository_UpdateById(t *testing.T) {
	const id = "0194f5c9-c3ca-7725-bbbc-4bafcd25c632"
	tcs := map[string]userTestCase{
		"SUCCESS": {
			reg: func(ctrl *gomock.Controller, conn *_dbMock.MockConn, _ *_uuidMock.MockGenerator) {
				result := _dbMock.NewMockResult(ctrl)
				result.EXPECT().RowsAffected().Times(1).Return(int64(1))

				conn.EXPECT().Execute(gomock.Any(), qUpdateUserById, gomock.Any(), gomock.Any(), gomock.Any(),
					gomock.Any(), gomock.Any()).Return(result, nil)
			},
			cmd: func(repo *UserRepository) error {
				return repo.UpdateById(context.Background(), &domain.User{ID: id})
			},
			exp: func(err error) {
				require.NoError(t, err)
			},
		},
		"FAILED user field must not be empty": {
			reg: func(ctrl *gomock.Controller, conn *_dbMock.MockConn, _ *_uuidMock.MockGenerator) {
				conn.EXPECT().Execute(gomock.Any(), qUpdateUserById, gomock.Any(), gomock.Any(), gomock.Any(),
					gomock.Any(), gomock.Any()).Return(nil, proxerr.New(database.ErrNotNullViolation, ""))
			},
			cmd: func(repo *UserRepository) error {
				return repo.UpdateById(context.Background(), &domain.User{ID: id})
			},
			exp: func(err error) {
				require.Equal(t, ErrUserFieldMustNotBeEmpty, err)
			},
		},
		"FAILED user not found": {
			reg: func(ctrl *gomock.Controller, conn *_dbMock.MockConn, _ *_uuidMock.MockGenerator) {
				conn.EXPECT().Execute(gomock.Any(), qUpdateUserById, gomock.Any(), gomock.Any(), gomock.Any(),
					gomock.Any(), gomock.Any()).Return(nil, proxerr.New(database.ErrForeignKeyViolation, ""))
			},
			cmd: func(repo *UserRepository) error {
				return repo.UpdateById(context.Background(), &domain.User{ID: id})
			},
			exp: func(err error) {
				require.Equal(t, ErrUserNotFound, err)
			},
		},
		"FAILED no rows were affected": {
			reg: func(ctrl *gomock.Controller, conn *_dbMock.MockConn, _ *_uuidMock.MockGenerator) {
				result := _dbMock.NewMockResult(ctrl)
				result.EXPECT().RowsAffected().Times(1).Return(int64(0))

				conn.EXPECT().Execute(gomock.Any(), qUpdateUserById, gomock.Any(), gomock.Any(), gomock.Any(),
					gomock.Any(), gomock.Any()).Return(result, nil)
			},
			cmd: func(repo *UserRepository) error {
				return repo.UpdateById(context.Background(), &domain.User{ID: id})
			},
			exp: func(err error) {
				require.Error(t, err)
			},
		},
		"FAILED": {
			reg: func(ctrl *gomock.Controller, conn *_dbMock.MockConn, _ *_uuidMock.MockGenerator) {
				conn.EXPECT().Execute(gomock.Any(), qUpdateUserById, gomock.Any(), gomock.Any(), gomock.Any(),
					gomock.Any(), gomock.Any()).Return(nil, errors.New(""))
			},
			cmd: func(repo *UserRepository) error {
				return repo.UpdateById(context.Background(), &domain.User{ID: id})
			},
			exp: func(err error) {
				require.Error(t, err)
			},
		},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			runUserTestCase(t, &tc)
		})
	}
}

func TestUserRepository_DeleteById(t *testing.T) {
	const id = "0194f5da-6239-70a8-b70d-bebbe9c835b6"
	tcs := map[string]userTestCase{
		"SUCCESS": {
			reg: func(ctrl *gomock.Controller, conn *_dbMock.MockConn, _ *_uuidMock.MockGenerator) {
				result := _dbMock.NewMockResult(ctrl)
				result.EXPECT().RowsAffected().Times(1).Return(int64(1))

				conn.EXPECT().Execute(gomock.Any(), qDeleteUserById, gomock.Any()).Times(1).Return(result, nil)
			},
			cmd: func(repo *UserRepository) error {
				return repo.DeleteById(context.Background(), id)
			},
			exp: func(err error) {
				require.NoError(t, err)
			},
		},
		"FAILED no rows were affected": {
			reg: func(ctrl *gomock.Controller, conn *_dbMock.MockConn, _ *_uuidMock.MockGenerator) {
				result := _dbMock.NewMockResult(ctrl)
				result.EXPECT().RowsAffected().Times(1).Return(int64(0))

				conn.EXPECT().Execute(gomock.Any(), qDeleteUserById, gomock.Any()).Times(1).Return(result, nil)
			},
			cmd: func(repo *UserRepository) error {
				return repo.DeleteById(context.Background(), id)
			},
			exp: func(err error) {
				require.Error(t, err)
			},
		},
		"FAILED": {
			reg: func(ctrl *gomock.Controller, conn *_dbMock.MockConn, _ *_uuidMock.MockGenerator) {
				conn.EXPECT().Execute(gomock.Any(), qDeleteUserById, gomock.Any()).Times(1).
					Return(nil, errors.New(""))
			},
			cmd: func(repo *UserRepository) error {
				return repo.DeleteById(context.Background(), id)
			},
			exp: func(err error) {
				require.Error(t, err)
			},
		},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			runUserTestCase(t, &tc)
		})
	}
}

// runUserTestCase should be called by [testing.T.Run]
func runUserTestCase(t *testing.T, tc *userTestCase) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	conn := _dbMock.NewMockConn(ctrl)
	idGen := _uuidMock.NewMockGenerator(ctrl)
	if tc.reg != nil {
		tc.reg(ctrl, conn, idGen)
	}

	logger := slog.NewLogger(stdslog.New(slogzap.Option{}.NewZapHandler()))
	repo := NewUserRepository(conn, logger, idGen)

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
