package postgres

import (
	"context"
	"github.com/adanyl0v/pocket-ideas/internal/repository"
	"github.com/adanyl0v/pocket-ideas/pkg/database"
	"github.com/adanyl0v/pocket-ideas/pkg/log"
)

type Repository struct {
	conn   database.Conn
	logger log.Logger
}

func (r *Repository) Begin(ctx context.Context) (repository.Tx, error) {
	tx, err := r.conn.Begin(ctx)
	if err != nil {
		r.logger.WithError(err).Error("failed to begin a user repository transaction")
		return nil, err
	}

	r.logger.Debug("begun a user repository transaction")
	return tx, err
}
