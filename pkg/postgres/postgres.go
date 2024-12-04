package postgres

import (
	"auth-service/internal/common"
	"auth-service/pkg/config"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

var (
	// create a new pool for every test
	//pool *pgxpool.Pool
	once sync.Once
)

const (
	maxConn         = 50
	maxConnLifetime = 3 * time.Minute
	minConns        = 10
)

func NewPG(config *config.Config) (*pgxpool.Pool, error) {
	//once.Do(func() {
	ctx := context.Background()
	dataSourceName := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=%s password=%s",
		config.DB.PostgresQL.PostgresqlHost,
		config.DB.PostgresQL.PostgresqlPort,
		config.DB.PostgresQL.PostgresqlUser,
		config.DB.PostgresQL.PostgresqlDatabase,
		config.DB.PostgresQL.PostgresqlSSLMode,
		config.DB.PostgresQL.PostgresqlPassword,
	)

	poolCfg, err := pgxpool.ParseConfig(dataSourceName)
	if err != nil {
		logrus.Error(err)
	}

	poolCfg.MaxConns = maxConn
	poolCfg.MaxConnLifetime = maxConnLifetime
	poolCfg.MinConns = minConns

	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		logrus.Error(err)
	}

	//})

	err = pool.Ping(ctx)
	if err != nil {
		pool.Close()
		return nil, errors.Wrap(common.ErrConnectionDB, err.Error())
	}

	return pool, nil
}
