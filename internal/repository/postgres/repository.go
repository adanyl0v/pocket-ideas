package postgres

import (
	pgdb "github.com/adanyl0v/pocket-ideas/pkg/database/postgres"
	"github.com/adanyl0v/pocket-ideas/pkg/log"
)

type repository struct {
	conn   pgdb.Conn
	logger log.Logger
}
