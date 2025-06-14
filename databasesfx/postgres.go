package databasesfx

import (
	"context"
	"fmt"
	"io/fs"
	"net"
	"strconv"
	"time"

	"github.com/golang-migrate/migrate/v4"
	migratepgx "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/fx"
)

type PostgresConfig struct {
	ApplicationName       string        `required:"true" mapstructure:"application_name"         yaml:"application_name"         json:"application_name"`
	Timezone              string        `required:"true" mapstructure:"timezone"                 yaml:"timezone"                 json:"timezone"                 default:"UTC"`
	DBName                string        `required:"true" mapstructure:"dbname"                   yaml:"dbname"                   json:"dbname"`
	Schema                string        `required:"true" mapstructure:"schema"                   yaml:"schema"                   json:"schema"                   default:"public"`
	Host                  string        `required:"true" mapstructure:"host"                     yaml:"host"                     json:"host"                     default:"localhost"`
	SslMode               string        `required:"true" mapstructure:"ssl_mode"                 yaml:"ssl_mode"                 json:"ssl_mode"                 default:"disable"`
	Password              string        `required:"true" mapstructure:"password"                 yaml:"password"                 json:"password"`
	Username              string        `required:"true" mapstructure:"username"                 yaml:"username"                 json:"username"`
	MaxIdleConnection     int           `required:"true" mapstructure:"max_idle_connections"     yaml:"max_idle_connections"     json:"max_idle_connections"`
	ConnectionTimeout     time.Duration `required:"true" mapstructure:"connection_timeout"       yaml:"connection_timeout"       json:"connection_timeout"       default:"5s"`
	MaxOpenConnections    int           `required:"true" mapstructure:"max_open_connections"     yaml:"max_open_connections"     json:"max_open_connections"`
	MaxConnectionLifetime time.Duration `required:"true" mapstructure:"max_connection_lifetime"  yaml:"max_connection_lifetime"  json:"max_connection_lifetime"`
	MaxConnectionIdleTime time.Duration `required:"true" mapstructure:"max_connection_idle_time" yaml:"max_connection_idle_time" json:"max_connection_idle_time"`
	Port                  uint16        `required:"true" mapstructure:"port"                     yaml:"port"                     json:"port"                     default:"5432"`
}

func NewPostgresMigrations(fs fs.FS, cfg PostgresConfig, migrations, migrationsTable string) (*migrate.Migrate, error) {
	sourceDriver, err := iofs.New(fs, migrations)
	if err != nil {
		return nil, err
	}

	pgxConfig, err := pgx.ParseConfig(cfg.ConnectionString())
	if err != nil {
		return nil, err
	}

	db, err := migratepgx.WithInstance(stdlib.OpenDB(*pgxConfig), &migratepgx.Config{
		MigrationsTable:       migrationsTable,
		DatabaseName:          cfg.DBName,
		SchemaName:            cfg.Schema,
		StatementTimeout:      60 * time.Second,
		MigrationsTableQuoted: false,
		MultiStatementEnabled: true,
	})
	if err != nil {
		return nil, err
	}

	return migrate.NewWithInstance("iofs", sourceDriver, "pgx5", db)
}

func PostgresModule(cfg PostgresConfig) fx.Option {
	return fx.Module("Databases-Postgres", fx.Provide(
		func(lc fx.Lifecycle) (*pgxpool.Pool, error) {
			ctx, cancel := context.WithTimeout(context.Background(), cfg.ConnectionTimeout)
			defer cancel()
			conn, err := pgxpool.New(ctx, cfg.ConnectionString())
			if err != nil {
				return nil, err
			}

			lc.Append(fx.StopHook(func(_ context.Context) error {
				conn.Close()
				return nil
			}))

			return conn, nil
		}),
	)
}

func PostgresMigrationsModule(mig fs.FS, cfg PostgresConfig, migrations string) fx.Option {
	return fx.Module("Databases-Postgres-Migrations",
		fx.Provide(func() (*migrate.Migrate, error) {
			return NewPostgresMigrations(mig, cfg, migrations, "migrations")
		}),
	)
}

func (p *PostgresConfig) ConnectionString() string {
	if p.Schema == "" {
		p.Schema = "public"
	}

	host := net.JoinHostPort(p.Host, strconv.FormatInt(int64(p.Port), 10))
	return fmt.Sprintf(
		"postgres://%s:%s@%s/%s?search_path=%s&sslmode=%s&application_name=%s&timezone=%s",
		p.Username,
		p.Password,
		host,
		p.DBName,
		p.Schema,
		p.SslMode,
		p.ApplicationName,
		p.Timezone,
	)
}
