package postgres

import (
	"context"
	"errors"
	_pgMock "github.com/adanyl0v/pocket-ideas/mocks/pkg/database/postgres/pgx"
	"github.com/adanyl0v/pocket-ideas/pkg/database"
	"github.com/adanyl0v/pocket-ideas/pkg/log/slog"
	"github.com/adanyl0v/pocket-ideas/pkg/proxerr"
	"github.com/golang/mock/gomock"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	slogzap "github.com/samber/slog-zap/v2"
	"github.com/stretchr/testify/require"
	stdslog "log/slog"
	"testing"
)

func TestRow_Scan(t *testing.T) {
	type rowTestCase struct {
		reg func(_ *gomock.Controller, row *_pgMock.MockDriverRow)
		cmd func(row *Row) error
		exp func(err error)
	}

	tcs := map[string]rowTestCase{
		"SUCCESS": {
			reg: func(_ *gomock.Controller, row *_pgMock.MockDriverRow) {
				row.EXPECT().Scan().Times(1).Return(nil)
			},
			cmd: func(row *Row) error {
				return row.Scan()
			},
			exp: func(err error) {
				require.NoError(t, err)
			},
		},
		"FAILED no rows in result set": {
			reg: func(_ *gomock.Controller, row *_pgMock.MockDriverRow) {
				row.EXPECT().Scan().Times(1).Return(pgx.ErrNoRows)
			},
			cmd: func(row *Row) error {
				return row.Scan()
			},
			exp: func(err error) {
				require.Equal(t, database.ErrNoRows, err)
			},
		},
		"FAILED": {
			reg: func(_ *gomock.Controller, row *_pgMock.MockDriverRow) {
				row.EXPECT().Scan().Times(1).Return(errors.New(""))
			},
			cmd: func(row *Row) error {
				return row.Scan()
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

			driverRow := _pgMock.NewMockDriverRow(ctrl)
			if tc.reg != nil {
				tc.reg(ctrl, driverRow)
			}

			var err error
			if tc.cmd != nil {
				err = tc.cmd(&Row{row: driverRow})
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
		})
	}
}

func TestRows_Scan(t *testing.T) {
	type rowsTestCase struct {
		reg func(_ *gomock.Controller, rows *_pgMock.MockDriverRows)
		cmd func(rows *Rows) error
		exp func(err error)
	}

	tcs := map[string]rowsTestCase{
		"SUCCESS": {
			reg: func(_ *gomock.Controller, rows *_pgMock.MockDriverRows) {
				rows.EXPECT().Scan().Times(1).Return(nil)
			},
			cmd: func(rows *Rows) error {
				return rows.Scan()
			},
			exp: func(err error) {
				require.NoError(t, err)
			},
		},
		"FAILED no rows in result set": {
			reg: func(_ *gomock.Controller, rows *_pgMock.MockDriverRows) {
				rows.EXPECT().Scan().Times(1).Return(pgx.ErrNoRows)
			},
			cmd: func(rows *Rows) error {
				return rows.Scan()
			},
			exp: func(err error) {
				require.Equal(t, database.ErrNoRows, err)
			},
		},
		"FAILED": {
			reg: func(_ *gomock.Controller, rows *_pgMock.MockDriverRows) {
				rows.EXPECT().Scan().Times(1).Return(errors.New(""))
			},
			cmd: func(rows *Rows) error {
				return rows.Scan()
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

			driverRows := _pgMock.NewMockDriverRows(ctrl)
			if tc.reg != nil {
				tc.reg(ctrl, driverRows)
			}

			var err error
			if tc.cmd != nil {
				err = tc.cmd(&Rows{rows: driverRows})
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
		})
	}
}

type (
	connTestCaseRegister func(_ *gomock.Controller, conn *_pgMock.MockDriverConn)
	connTestCaseCommand  func(conn *Conn) error
	connTestCaseExpect   func(err error)

	connTestCase struct {
		reg connTestCaseRegister
		cmd connTestCaseCommand
		exp connTestCaseExpect
	}
)

func TestConn_Execute(t *testing.T) {
	tcs := map[string]connTestCase{
		"SUCCESS": {
			reg: func(_ *gomock.Controller, conn *_pgMock.MockDriverConn) {
				conn.EXPECT().Exec(gomock.Any(), "").Return(pgconn.CommandTag{}, nil)
			},
			cmd: func(conn *Conn) error {
				_, err := conn.Execute(context.Background(), "")
				return err
			},
			exp: func(err error) {
				require.NoError(t, err)
			},
		},
		"FAILED check violation": {
			reg: func(_ *gomock.Controller, conn *_pgMock.MockDriverConn) {
				conn.EXPECT().Exec(gomock.Any(), "").
					Return(pgconn.CommandTag{}, &pgconn.PgError{Code: pgerrcode.CheckViolation})
			},
			cmd: func(conn *Conn) error {
				_, err := conn.Execute(context.Background(), "")
				return err
			},
			exp: func(err error) {
				require.Equal(t, database.ErrCheckViolation, err)
			},
		},
		"FAILED unique violation": {
			reg: func(_ *gomock.Controller, conn *_pgMock.MockDriverConn) {
				conn.EXPECT().Exec(gomock.Any(), "").
					Return(pgconn.CommandTag{}, &pgconn.PgError{Code: pgerrcode.UniqueViolation})
			},
			cmd: func(conn *Conn) error {
				_, err := conn.Execute(context.Background(), "")
				return err
			},
			exp: func(err error) {
				require.Equal(t, database.ErrUniqueViolation, err)
			},
		},
		"FAILED not null violation": {
			reg: func(_ *gomock.Controller, conn *_pgMock.MockDriverConn) {
				conn.EXPECT().Exec(gomock.Any(), "").
					Return(pgconn.CommandTag{}, &pgconn.PgError{Code: pgerrcode.NotNullViolation})
			},
			cmd: func(conn *Conn) error {
				_, err := conn.Execute(context.Background(), "")
				return err
			},
			exp: func(err error) {
				require.Equal(t, database.ErrNotNullViolation, err)
			},
		},
		"FAILED foreign key violation": {
			reg: func(_ *gomock.Controller, conn *_pgMock.MockDriverConn) {
				conn.EXPECT().Exec(gomock.Any(), "").
					Return(pgconn.CommandTag{}, &pgconn.PgError{Code: pgerrcode.ForeignKeyViolation})
			},
			cmd: func(conn *Conn) error {
				_, err := conn.Execute(context.Background(), "")
				return err
			},
			exp: func(err error) {
				require.Equal(t, database.ErrForeignKeyViolation, err)
			},
		},
		"FAILED": {
			reg: func(_ *gomock.Controller, conn *_pgMock.MockDriverConn) {
				conn.EXPECT().Exec(gomock.Any(), "").
					Return(pgconn.CommandTag{}, errors.New(""))
			},
			cmd: func(conn *Conn) error {
				_, err := conn.Execute(context.Background(), "")
				return err
			},
			exp: func(err error) {
				require.Error(t, err)
			},
		},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			runConnTestCase(t, &tc)
		})
	}
}

func TestConn_Query(t *testing.T) {
	tcs := map[string]connTestCase{
		"SUCCESS": {
			reg: func(ctrl *gomock.Controller, conn *_pgMock.MockDriverConn) {
				conn.EXPECT().Query(gomock.Any(), "").Return(nil, nil)
			},
			cmd: func(conn *Conn) error {
				_, err := conn.Query(context.Background(), "")
				return err
			},
			exp: func(err error) {
				require.NoError(t, err)
			},
		},
		"FAILED no rows in result set": {
			reg: func(ctrl *gomock.Controller, conn *_pgMock.MockDriverConn) {
				conn.EXPECT().Query(gomock.Any(), "").Return(nil, pgx.ErrNoRows)
			},
			cmd: func(conn *Conn) error {
				_, err := conn.Query(context.Background(), "")
				return err
			},
			exp: func(err error) {
				require.Equal(t, database.ErrNoRows, err)
			},
		},
		"FAILED": {
			reg: func(ctrl *gomock.Controller, conn *_pgMock.MockDriverConn) {
				conn.EXPECT().Query(gomock.Any(), "").Return(nil, errors.New(""))
			},
			cmd: func(conn *Conn) error {
				_, err := conn.Query(context.Background(), "")
				return err
			},
			exp: func(err error) {
				require.Error(t, err)
			},
		},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			runConnTestCase(t, &tc)
		})
	}
}

func TestConn_QueryRow(t *testing.T) {
	tcs := map[string]connTestCase{
		"SUCCESS": {
			reg: func(_ *gomock.Controller, conn *_pgMock.MockDriverConn) {
				conn.EXPECT().QueryRow(gomock.Any(), "").Return(nil)
			},
			cmd: func(conn *Conn) error {
				_ = conn.QueryRow(context.Background(), "")
				return nil
			},
			exp: func(err error) {
				require.NoError(t, err)
			},
		},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			runConnTestCase(t, &tc)
		})
	}
}

func TestConn_Begin(t *testing.T) {
	tcs := map[string]connTestCase{
		"SUCCESS": {
			reg: func(_ *gomock.Controller, conn *_pgMock.MockDriverConn) {
				conn.EXPECT().Begin(gomock.Any()).Return(nil, nil)
			},
			cmd: func(conn *Conn) error {
				_, err := conn.Begin(context.Background())
				return err
			},
			exp: func(err error) {
				require.NoError(t, err)
			},
		},
		"FAILED": {
			reg: func(_ *gomock.Controller, conn *_pgMock.MockDriverConn) {
				conn.EXPECT().Begin(gomock.Any()).Return(nil, errors.New(""))
			},
			cmd: func(conn *Conn) error {
				_, err := conn.Begin(context.Background())
				return err
			},
			exp: func(err error) {
				require.Error(t, err)
			},
		},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			runConnTestCase(t, &tc)
		})
	}
}

func runConnTestCase(t *testing.T, tc *connTestCase) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	driverConn := _pgMock.NewMockDriverConn(ctrl)
	if tc.reg != nil {
		tc.reg(ctrl, driverConn)
	}

	var err error
	if tc.cmd != nil {
		err = tc.cmd(&Conn{
			conn:   driverConn,
			logger: slog.NewLogger(stdslog.New(slogzap.Option{}.NewZapHandler())),
		})
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

type (
	txTestCaseRegister func(_ *gomock.Controller, tx *_pgMock.MockDriverTx)
	txTestCaseCommand  func(tx *Tx) error
	txTestCaseExpect   func(err error)

	txTestCase struct {
		reg txTestCaseRegister
		cmd txTestCaseCommand
		exp txTestCaseExpect
	}
)

func TestTx_Commit(t *testing.T) {
	tcs := map[string]txTestCase{
		"SUCCESS": {
			reg: func(_ *gomock.Controller, tx *_pgMock.MockDriverTx) {
				tx.EXPECT().Commit(gomock.Any()).Return(nil)
			},
			cmd: func(conn *Tx) error {
				return conn.Commit(context.Background())
			},
			exp: func(err error) {
				require.NoError(t, err)
			},
		},
		"FAILED not a transaction": {
			cmd: func(conn *Tx) error {
				conn.conn = nil
				return conn.Commit(context.Background())
			},
			exp: func(err error) {
				require.Equal(t, ErrNotTransaction, err)
			},
		},
		"FAILED": {
			reg: func(_ *gomock.Controller, tx *_pgMock.MockDriverTx) {
				tx.EXPECT().Commit(gomock.Any()).Return(errors.New(""))
			},
			cmd: func(conn *Tx) error {
				return conn.Commit(context.Background())
			},
			exp: func(err error) {
				require.Error(t, err)
			},
		},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			runTxTestCase(t, &tc)
		})
	}
}

func TestTx_Rollback(t *testing.T) {
	tcs := map[string]txTestCase{
		"SUCCESS": {
			reg: func(_ *gomock.Controller, tx *_pgMock.MockDriverTx) {
				tx.EXPECT().Rollback(gomock.Any()).Return(nil)
			},
			cmd: func(conn *Tx) error {
				return conn.Rollback(context.Background())
			},
			exp: func(err error) {
				require.NoError(t, err)
			},
		},
		"FAILED not a transaction": {
			cmd: func(conn *Tx) error {
				conn.conn = nil
				return conn.Rollback(context.Background())
			},
			exp: func(err error) {
				require.Equal(t, ErrNotTransaction, err)
			},
		},
		"FAILED": {
			reg: func(_ *gomock.Controller, tx *_pgMock.MockDriverTx) {
				tx.EXPECT().Rollback(gomock.Any()).Return(errors.New(""))
			},
			cmd: func(conn *Tx) error {
				return conn.Rollback(context.Background())
			},
			exp: func(err error) {
				require.Error(t, err)
			},
		},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			runTxTestCase(t, &tc)
		})
	}
}

func runTxTestCase(t *testing.T, tc *txTestCase) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	driverTx := _pgMock.NewMockDriverTx(ctrl)
	if tc.reg != nil {
		tc.reg(ctrl, driverTx)
	}

	var err error
	if tc.cmd != nil {
		err = tc.cmd(&Tx{
			Conn{
				conn:   driverTx,
				logger: slog.NewLogger(stdslog.New(slogzap.Option{}.NewZapHandler())),
			},
		})
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
