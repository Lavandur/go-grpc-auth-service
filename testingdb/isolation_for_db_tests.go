package testingdb

import (
	"auth-service/pkg/config"
	"auth-service/pkg/postgres"
	"context"
	"encoding/base64"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
	"math/rand"
	"os"
	"strings"
	"sync"
	"time"
)

type TestingT interface {
	require.TestingT

	Log(args ...any)
	Logf(format string, args ...any)
	Cleanup(f func())
}

func NewWithIsolatedDatabase(t TestingT) *Postgres {
	return newPostgres(t).cloneFromReference()
}

func prepareConfig() *config.Config {
	return &config.Config{
		App: config.App{
			DB: config.DB{
				PostgresQL: config.PostgresQL{
					PostgresqlHost:     "localhost",
					PostgresqlPort:     "5432",
					PostgresqlUser:     "auth",
					PostgresqlPassword: "auth",
					PostgresqlDatabase: "auth",
					PostgresqlSSLMode:  "disable",
				},
			},
		},
	}
}

type Postgres struct {
	t TestingT

	conf *config.Config
	ref  string

	pool     *pgxpool.Pool
	poolOnce sync.Once
}

func newPostgres(t TestingT) *Postgres {
	conf := prepareConfig()

	refDB := os.Getenv("REF_DB_NAME")
	if refDB == "" {
		refDB = "reference"
	}

	return &Postgres{
		t: t,

		conf: conf,
		ref:  refDB,
	}
}

func (p *Postgres) DB() *pgxpool.Pool {
	p.poolOnce.Do(func() {
		p.pool = open(p.t, p.conf)
	})

	return p.pool
}

func (p *Postgres) cloneFromReference() *Postgres {
	newDBName := getUniqueDBName(p.ref, p.t)

	p.t.Log("database name for this test:", newDBName)

	sql := fmt.Sprintf(
		`CREATE DATABASE %q WITH TEMPLATE %q;`,
		newDBName,
		p.ref,
	)

	_, err := p.DB().Exec(context.Background(), sql)
	require.NoError(p.t, err)

	p.t.Cleanup(func() {
		sql := fmt.Sprintf(`DROP DATABASE %q WITH (FORCE);`, newDBName)

		ctx, done := context.WithTimeout(context.Background(), time.Minute)
		defer done()

		_, err := p.DB().Exec(ctx, sql)
		require.NoError(p.t, err)
	})

	return p.replaceDBName(newDBName)
}

func (p *Postgres) replaceDBName(newDBName string) *Postgres {
	o := p.clone()
	p.conf.DB.PostgresqlDatabase = newDBName
	o.conf = p.conf

	return o
}

func (p *Postgres) clone() *Postgres {
	return &Postgres{
		t: p.t,

		conf: p.conf,
		ref:  p.ref,
	}
}

func getUniqueDBName(prefix string, t TestingT) string {
	dbName := strings.Builder{}
	dbName.WriteString(prefix)
	dbName.WriteRune('-')

	bs := make([]byte, 6)

	_, err := rand.Read(bs)
	require.NoError(t, err)

	r := base64.RawURLEncoding.EncodeToString(bs)
	dbName.WriteString(r)

	return dbName.String()
}

func open(t TestingT, conf *config.Config) *pgxpool.Pool {

	pg, err := postgres.NewPG(conf)
	require.NoError(t, err)

	// Close connection after the test is completed.
	t.Cleanup(func() {
		pg.Close()
	})

	return pg
}
