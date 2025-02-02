package postgres

import (
	"github.com/adanyl0v/pocket-ideas/pkg/log"
	"github.com/adanyl0v/pocket-ideas/pkg/log/slog"
	slogzap "github.com/samber/slog-zap/v2"
	stdslog "log/slog"
	"os"
	"testing"
)

var testCaseEmptyLogger log.Logger

func TestMain(m *testing.M) {
	testCaseEmptyLogger = slog.NewLogger(stdslog.New(slogzap.Option{}.NewZapHandler()))
	os.Exit(m.Run())
}
