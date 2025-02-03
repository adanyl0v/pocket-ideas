package postgres

import (
	"context"
	"errors"
	"github.com/adanyl0v/pocket-ideas/internal/domain"
	"github.com/adanyl0v/pocket-ideas/internal/repository"
	_pgdbMock "github.com/adanyl0v/pocket-ideas/mocks/pkg/database/postgres"
	_uuidMock "github.com/adanyl0v/pocket-ideas/mocks/pkg/uuid"
	pgdb "github.com/adanyl0v/pocket-ideas/pkg/database/postgres"
	"github.com/adanyl0v/pocket-ideas/pkg/proxerr"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"testing"
)

type (
	userTestCaseRegister     func(ctrl *gomock.Controller, conn *_pgdbMock.MockConn, idGen *_uuidMock.MockGenerator)
	userTestCaseErrorChecker func(err error)
	userTestCaseRepoReceiver func(repo *UserRepository) error

	userTestCase struct {
		name     string
		register userTestCaseRegister
		receiver userTestCaseRepoReceiver
		checker  userTestCaseErrorChecker
	}
)

func TestNewUserRepository(t *testing.T) {
	NewUserRepository(nil, nil, nil)
}

func TestUserRepository_Begin(t *testing.T) {
	testCases := []userTestCase{
		{
			name: testCaseSuccess,
			register: func(ctrl *gomock.Controller, conn *_pgdbMock.MockConn, _ *_uuidMock.MockGenerator) {
				tx := _pgdbMock.NewMockTx(ctrl)
				conn.EXPECT().Begin(gomock.Any()).Times(1).Return(tx, nil)
			},
			receiver: func(repo *UserRepository) error {
				_, err := repo.Begin(context.Background())
				return err
			},
			checker: func(err error) {
				require.NoError(t, err)
			},
		},
		{
			name: testCaseFailure,
			register: func(_ *gomock.Controller, conn *_pgdbMock.MockConn, _ *_uuidMock.MockGenerator) {
				conn.EXPECT().Begin(gomock.Any()).Times(1).Return(nil, errors.New(""))
			},
			receiver: func(repo *UserRepository) error {
				_, err := repo.Begin(context.Background())
				return err
			},
			checker: func(err error) {
				require.NotNil(t, err)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			runUserTestCase(t, tc)
		})
	}
}

func TestUserRepository_WithTx(t *testing.T) {
	testCases := []userTestCase{
		{
			name: testCaseSuccess,
			receiver: func(repo *UserRepository) error {
				r := repo.WithTx(new(_pgdbMock.MockTx))
				require.NotNil(t, r)
				return nil
			},
		},
		{
			name: testCaseFailure,
			receiver: func(repo *UserRepository) error {
				r := repo.WithTx(repository.Tx(nil))
				require.Nil(t, r)
				return nil
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			runUserTestCase(t, tc)
		})
	}
}

func TestUserRepository_Create(t *testing.T) {
	testCases := []userTestCase{
		{
			name: testCaseSuccess,
			register: func(_ *gomock.Controller, conn *_pgdbMock.MockConn, idGen *_uuidMock.MockGenerator) {
				idGen.EXPECT().NewV7().Times(1).Return("", nil)
				conn.EXPECT().Exec(gomock.Any(), insertUserQuery, gomock.Any(), gomock.Any(), gomock.Any(),
					gomock.Any(), gomock.Any(), gomock.Any()).Times(1).Return(nil)
			},
			receiver: func(repo *UserRepository) error {
				return repo.Create(context.Background(), new(domain.User))
			},
			checker: func(err error) {
				require.Nil(t, err)
			},
		},
		{
			name: msgFailedToGenerateUserID,
			register: func(_ *gomock.Controller, _ *_pgdbMock.MockConn, idGen *_uuidMock.MockGenerator) {
				idGen.EXPECT().NewV7().Times(1).Return("", errGenerateUserId)
			},
			receiver: func(repo *UserRepository) error {
				return repo.Create(context.Background(), new(domain.User))
			},
			checker: func(err error) {
				require.Equal(t, errGenerateUserId, err)
			},
		},
		{
			name: msgUserAlreadyExists,
			register: func(_ *gomock.Controller, conn *_pgdbMock.MockConn, idGen *_uuidMock.MockGenerator) {
				idGen.EXPECT().NewV7().Times(1).Return("", nil)
				conn.EXPECT().Exec(gomock.Any(), insertUserQuery, gomock.Any(), gomock.Any(), gomock.Any(),
					gomock.Any(), gomock.Any(), gomock.Any()).Times(1).
					Return(proxerr.New(pgdb.ErrUniqueViolation, ""))
			},
			receiver: func(repo *UserRepository) error {
				return repo.Create(context.Background(), new(domain.User))
			},
			checker: func(err error) {
				require.Equal(t, ErrUserExists, err)
			},
		},
		{
			name: testCaseFailure,
			register: func(ctrl *gomock.Controller, conn *_pgdbMock.MockConn, idGen *_uuidMock.MockGenerator) {
				idGen.EXPECT().NewV7().Times(1).Return("", nil)
				conn.EXPECT().Exec(gomock.Any(), insertUserQuery, gomock.Any(), gomock.Any(), gomock.Any(),
					gomock.Any(), gomock.Any(), gomock.Any()).Times(1).
					Return(errors.New(""))
			},
			receiver: func(repo *UserRepository) error {
				return repo.Create(context.Background(), new(domain.User))
			},
			checker: func(err error) {
				require.NotNil(t, err)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			runUserTestCase(t, tc)
		})
	}
}

func TestUserRepository_GetByID(t *testing.T) {
	testCases := []userTestCase{
		{
			name: testCaseSuccess,
			register: func(ctrl *gomock.Controller, conn *_pgdbMock.MockConn, _ *_uuidMock.MockGenerator) {
				row := _pgdbMock.NewMockRow(ctrl)
				row.EXPECT().Scan(gomock.Any(), gomock.Any(), gomock.Any(),
					gomock.Any(), gomock.Any()).Times(1).Return(nil)
				conn.EXPECT().QueryRow(gomock.Any(), getUserByIdQuery, gomock.Any()).Times(1).Return(row)
			},
			receiver: func(repo *UserRepository) error {
				_, err := repo.GetByID(context.Background(), "")
				return err
			},
			checker: func(err error) {
				require.Nil(t, err)
			},
		},
		{
			name: msgUserNotFound,
			register: func(ctrl *gomock.Controller, conn *_pgdbMock.MockConn, idGen *_uuidMock.MockGenerator) {
				row := _pgdbMock.NewMockRow(ctrl)
				row.EXPECT().Scan(gomock.Any(), gomock.Any(), gomock.Any(),
					gomock.Any(), gomock.Any()).Times(1).Return(proxerr.New(pgdb.ErrNoRows, ""))
				conn.EXPECT().QueryRow(gomock.Any(), getUserByIdQuery, gomock.Any()).Times(1).Return(row)
			},
			receiver: func(repo *UserRepository) error {
				_, err := repo.GetByID(context.Background(), "")
				return err
			},
			checker: func(err error) {
				require.Equal(t, ErrUserNotFound, err)
			},
		},
		{
			name: testCaseFailure,
			register: func(ctrl *gomock.Controller, conn *_pgdbMock.MockConn, _ *_uuidMock.MockGenerator) {
				row := _pgdbMock.NewMockRow(ctrl)
				row.EXPECT().Scan(gomock.Any(), gomock.Any(), gomock.Any(),
					gomock.Any(), gomock.Any()).Times(1).Return(errors.New(""))
				conn.EXPECT().QueryRow(gomock.Any(), getUserByIdQuery, gomock.Any()).Times(1).Return(row)
			},
			receiver: func(repo *UserRepository) error {
				_, err := repo.GetByID(context.Background(), "")
				return err
			},
			checker: func(err error) {
				require.NotNil(t, err)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			runUserTestCase(t, tc)
		})
	}
}

func TestUserRepository_GetByEmail(t *testing.T) {
	testCases := []userTestCase{
		{
			name: testCaseSuccess,
			register: func(ctrl *gomock.Controller, conn *_pgdbMock.MockConn, _ *_uuidMock.MockGenerator) {
				row := _pgdbMock.NewMockRow(ctrl)
				row.EXPECT().Scan(gomock.Any(), gomock.Any(), gomock.Any(),
					gomock.Any(), gomock.Any()).Times(1).Return(nil)
				conn.EXPECT().QueryRow(gomock.Any(), getUserByEmailQuery, gomock.Any()).Times(1).Return(row)
			},
			receiver: func(repo *UserRepository) error {
				_, err := repo.GetByEmail(context.Background(), "")
				return err
			},
			checker: func(err error) {
				require.Nil(t, err)
			},
		},
		{
			name: msgUserNotFound,
			register: func(ctrl *gomock.Controller, conn *_pgdbMock.MockConn, idGen *_uuidMock.MockGenerator) {
				row := _pgdbMock.NewMockRow(ctrl)
				row.EXPECT().Scan(gomock.Any(), gomock.Any(), gomock.Any(),
					gomock.Any(), gomock.Any()).Times(1).Return(proxerr.New(pgdb.ErrNoRows, ""))
				conn.EXPECT().QueryRow(gomock.Any(), getUserByEmailQuery, gomock.Any()).Times(1).Return(row)
			},
			receiver: func(repo *UserRepository) error {
				_, err := repo.GetByEmail(context.Background(), "")
				return err
			},
			checker: func(err error) {
				require.Equal(t, ErrUserNotFound, err)
			},
		},
		{
			name: testCaseFailure,
			register: func(ctrl *gomock.Controller, conn *_pgdbMock.MockConn, _ *_uuidMock.MockGenerator) {
				row := _pgdbMock.NewMockRow(ctrl)
				row.EXPECT().Scan(gomock.Any(), gomock.Any(), gomock.Any(),
					gomock.Any(), gomock.Any()).Times(1).Return(errors.New(""))
				conn.EXPECT().QueryRow(gomock.Any(), getUserByEmailQuery, gomock.Any()).Times(1).Return(row)
			},
			receiver: func(repo *UserRepository) error {
				_, err := repo.GetByEmail(context.Background(), "")
				return err
			},
			checker: func(err error) {
				require.NotNil(t, err)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			runUserTestCase(t, tc)
		})
	}
}

func TestUserRepository_SelectAll(t *testing.T) {
	testCases := []userTestCase{
		{
			name: testCaseSuccess,
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
			receiver: func(repo *UserRepository) error {
				_, err := repo.SelectAll(context.Background())
				return err
			},
			checker: func(err error) {
				require.Nil(t, err)
			},
		},
		{
			name: msgFailedToSelectAllUsers,
			register: func(_ *gomock.Controller, conn *_pgdbMock.MockConn, _ *_uuidMock.MockGenerator) {
				conn.EXPECT().Query(gomock.Any(), selectAllUsersQuery).Times(1).Return(nil, errors.New(""))
			},
			receiver: func(repo *UserRepository) error {
				_, err := repo.SelectAll(context.Background())
				return err
			},
			checker: func(err error) {
				require.NotNil(t, err)
			},
		},
		{
			name: msgNoUsersFound,
			register: func(ctrl *gomock.Controller, conn *_pgdbMock.MockConn, _ *_uuidMock.MockGenerator) {
				rows := _pgdbMock.NewMockRows(ctrl)
				rows.EXPECT().Close().Times(1)
				rows.EXPECT().Next().Times(1).Return(false)
				rows.EXPECT().Err().Times(1).Return(nil)
				conn.EXPECT().Query(gomock.Any(), selectAllUsersQuery).Times(1).Return(rows, nil)
			},
			receiver: func(repo *UserRepository) error {
				_, err := repo.SelectAll(context.Background())
				return err
			},
			checker: func(err error) {
				require.Equal(t, ErrNoUsersFound, err)
			},
		},
		{
			name: msgNoUserRowsAffected,
			register: func(ctrl *gomock.Controller, conn *_pgdbMock.MockConn, _ *_uuidMock.MockGenerator) {
				rows := _pgdbMock.NewMockRows(ctrl)
				rows.EXPECT().Close().Times(1)
				rows.EXPECT().Next().Times(1).Return(false)
				rows.EXPECT().Err().Times(1).Return(errors.New(""))
				conn.EXPECT().Query(gomock.Any(), selectAllUsersQuery).Times(1).Return(rows, nil)
			},
			receiver: func(repo *UserRepository) error {
				_, err := repo.SelectAll(context.Background())
				return err
			},
			checker: func(err error) {
				require.NotNil(t, err)
			},
		},
		{
			name: msgFailedToScanAllUsers,
			register: func(ctrl *gomock.Controller, conn *_pgdbMock.MockConn, _ *_uuidMock.MockGenerator) {
				rows := _pgdbMock.NewMockRows(ctrl)
				rows.EXPECT().Close().Times(1)
				rows.EXPECT().Next().Times(1).Return(true)
				rows.EXPECT().Scan(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
					gomock.Any(), gomock.Any()).Times(1).Return(errors.New(""))
				conn.EXPECT().Query(gomock.Any(), selectAllUsersQuery).Times(1).Return(rows, nil)
			},
			receiver: func(repo *UserRepository) error {
				_, err := repo.SelectAll(context.Background())
				return err
			},
			checker: func(err error) {
				require.NotNil(t, err)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			runUserTestCase(t, tc)
		})
	}
}

// runUserTestCase should be called by [t.Run] to enable
// the Run gutter icon (at least in Goland from JetBrains)
//
// [t.Run]: https://pkg.go.dev/testing#T.Run
func runUserTestCase(t *testing.T, tc userTestCase) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	conn := _pgdbMock.NewMockConn(ctrl)
	idGen := _uuidMock.NewMockGenerator(ctrl)

	if tc.register != nil {
		tc.register(ctrl, conn, idGen)
	}

	repo := NewUserRepository(conn, testCaseEmptyLogger, idGen)

	var err error
	if tc.receiver != nil {
		err = tc.receiver(repo)
	}

	if err != nil {
		var pxErr proxerr.Error
		if errors.As(err, &pxErr) {
			err = pxErr.Unwrap()
		}
	}

	if tc.checker != nil {
		tc.checker(err)
	}
}
