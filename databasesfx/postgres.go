package databasesfx

import (
	"context"
	"fmt"
	"io/fs"
	"net"
	"strconv"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/fx"
)

type PostgresConfig struct {
	ApplicationName       string        `required:"true" mapstructure:"application_name" yaml:"application_name" json:"application_name"`
	Timezone              string        `required:"true" default:"UTC" mapstructure:"timezone" yaml:"timezone" json:"timezone"`
	DBName                string        `required:"true" mapstructure:"dbname" yaml:"dbname" json:"dbname"`
	Schema                string        `required:"true" default:"public" mapstructure:"schema" yaml:"schema" json:"schema"`
	Host                  string        `required:"true" default:"localhost" mapstructure:"host" yaml:"host" json:"host"`
	SslMode               string        `required:"true" default:"disable" mapstructure:"ssl_mode" yaml:"ssl_mode" json:"ssl_mode"`
	Password              string        `required:"true" yaml:"password" mapstructure:"password" json:"password"`
	Username              string        `required:"true" yaml:"username" mapstructure:"username" json:"username"`
	MaxIdleConnection     int           `required:"true" mapstructure:"max_idle_connections" yaml:"max_idle_connections" json:"max_idle_connections"`
	ConnectionTimeout     time.Duration `required:"true" default:"5s"  mapstructure:"connection_timeout" yaml:"connection_timeout" json:"connection_timeout"`
	MaxOpenConnections    int           `required:"true" mapstructure:"max_open_connections" yaml:"max_open_connections" json:"max_open_connections"`
	MaxConnectionLifetime time.Duration `required:"true" mapstructure:"max_connection_lifetime" yaml:"max_connection_lifetime" json:"max_connection_lifetime"`
	MaxConnectionIdleTime time.Duration `required:"true" mapstructure:"max_connection_idle_time" yaml:"max_connection_idle_time" json:"max_connection_idle_time"`
	Port                  uint16        `required:"true" default:"5432" mapstructure:"port" yaml:"port" json:"port"`
}

func New(fs fs.FS, cfg PostgresConfig, migrations string) (*migrate.Migrate, error) {
	sourceDriver, err := iofs.New(fs, migrations)
	if err != nil {
		return nil, err
	}

	return migrate.NewWithSourceInstance("iofs", sourceDriver, cfg.MigrationConnectionString())
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
			return New(mig, cfg, migrations)
		}),
	)
}

func (p *PostgresConfig) MigrationConnectionString() string {
	host := net.JoinHostPort(p.Host, strconv.FormatInt(int64(p.Port), 10))
	return fmt.Sprintf(
		"postgres://%s:%s@%s/%s?search_path=%s&sslmode=%s&application_name=%s&x-migrations-table=migrations&x-multi-statements=true",
		p.Username,
		p.Password,
		host,
		p.DBName,
		p.Schema,
		p.SslMode,
		p.ApplicationName,
	)
}

func (p *PostgresConfig) ConnectionString() string {
	if p.Schema == "" {
		return fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s application_name=%s",
			p.Host,
			p.Username,
			p.Password,
			p.DBName,
			p.Port,
			p.SslMode,
			p.Timezone,
			p.ApplicationName,
		)
	}

	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s search_path=%s application_name=%s",
		p.Host,
		p.Username,
		p.Password,
		p.DBName,
		p.Port,
		p.SslMode,
		p.Timezone,
		p.Schema,
		p.ApplicationName,
	)
}
