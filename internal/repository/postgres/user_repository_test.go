package postgres

import (
	"context"
	"errors"
	"github.com/adanyl0v/pocket-ideas/internal/domain"
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

func TestUserRepository_Create(t *testing.T) {
	testCases := []userTestCase{
		{
			name: "success",
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
			name: "failed to generate user id",
			register: func(_ *gomock.Controller, _ *_pgdbMock.MockConn, idGen *_uuidMock.MockGenerator) {
				idGen.EXPECT().NewV7().Times(1).Return("", errGenUserID)
			},
			receiver: func(repo *UserRepository) error {
				return repo.Create(context.Background(), new(domain.User))
			},
			checker: func(err error) {
				require.Equal(t, err, errGenUserID)
			},
		},
		{
			name: "user already exists",
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
				require.Equal(t, err, ErrUserExists)
			},
		},
		{
			name: "failure",
			register: func(_ *gomock.Controller, conn *_pgdbMock.MockConn, idGen *_uuidMock.MockGenerator) {
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

	// The [runUserTestCase] function should be called by [testing.T.Run]
	// to enable the Run gutter icon (at least in Goland from JetBrains)
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			runUserTestCase(t, tc)
		})
	}
}

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
