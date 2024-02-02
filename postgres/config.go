package postgres

import (
	"time"

	"github.com/rohanraj7316/utils/env"
	"github.com/rohanraj7316/utils/typecasting"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

const (
	ENV_POSTGRESQL_HOST               = "POSTGRESQL_HOST"
	ENV_POSTGRESQL_USER               = "POSTGRESQL_USER"
	ENV_POSTGRESQL_PORT               = "POSTGRESQL_PORT"
	ENV_POSTGRESQL_DB_NAME            = "POSTGRESQL_DB_NAME"
	ENV_POSTGRESQL_PASSWORD           = "POSTGRESQL_PASSWORD"
	ENV_POSTGRESQL_MAX_IDLE_CONN      = "POSTGRESQL_MAX_IDLE_CONN"
	ENV_POSTGRESQL_MAX_OPEN_CONN      = "POSTGRESQL_MAX_OPEN_CONN"
	ENV_POSTGRESQL_CONN_MAX_IDLE_TIME = "POSTGRESQL_CONN_MAX_IDLE_TIME"
	ENV_POSTGRESQL_CONN_MAX_LIFE_TIME = "POSTGRESQL_CONN_MAX_LIFE_TIME"

	CONSTANTS_SLOW_QUERY_THRESHOLD = 100 * time.Millisecond
)

type Config struct {
	User                            string
	Host                            string
	Port                            string
	DbName                          string
	Password                        string
	SslMode                         string
	ConnectionTimeout               int
	StatementTimeout                int
	IdleInTransactionSessionTimeout int
	MaxOpenConns                    int
	MaxIdleConns                    int
	ConnMaxLifetime                 time.Duration
	ConnMaxIdleTime                 time.Duration

	gormCfg *gorm.Config
}

var ConfigDefault = Config{
	SslMode: "disable",
	gormCfg: &gorm.Config{
		Logger: &log{
			SlowQueryThreshold: CONSTANTS_SLOW_QUERY_THRESHOLD,
			LogLevel:           logger.Error,
		},
		NamingStrategy: schema.NamingStrategy{
			SingularTable: false,
			NoLowerCase:   false,
		},
		DisableAutomaticPing: false,
	},
	MaxOpenConns:    100,
	MaxIdleConns:    10,
	ConnMaxLifetime: 4 * time.Second,
	ConnMaxIdleTime: 4 * time.Second,
}

func configDefault(config ...Config) (cfg Config, err error) {
	rConfigs := []string{
		ENV_POSTGRESQL_HOST,
		ENV_POSTGRESQL_USER,
		ENV_POSTGRESQL_DB_NAME,
		ENV_POSTGRESQL_PASSWORD,
		ENV_POSTGRESQL_PORT,
	}

	uConfigs := []string{
		ENV_POSTGRESQL_MAX_IDLE_CONN,
		ENV_POSTGRESQL_MAX_OPEN_CONN,
		ENV_POSTGRESQL_CONN_MAX_IDLE_TIME,
		ENV_POSTGRESQL_CONN_MAX_LIFE_TIME,
	}

	envCfg := env.EnvData(rConfigs, uConfigs)

	ConfigDefault.User = envCfg[ENV_POSTGRESQL_USER]
	ConfigDefault.Host = envCfg[ENV_POSTGRESQL_HOST]
	ConfigDefault.Port = envCfg[ENV_POSTGRESQL_PORT]
	ConfigDefault.DbName = envCfg[ENV_POSTGRESQL_DB_NAME]
	ConfigDefault.Password = envCfg[ENV_POSTGRESQL_PASSWORD]

	if envCfg[ENV_POSTGRESQL_MAX_OPEN_CONN] != "" {
		maxOpenConns, err := typecasting.ConvertStringToInt(envCfg[ENV_POSTGRESQL_MAX_OPEN_CONN],
			ENV_POSTGRESQL_MAX_OPEN_CONN)
		if err != nil {
			return cfg, err
		} else {
			ConfigDefault.MaxOpenConns = maxOpenConns
		}
	}

	if envCfg[ENV_POSTGRESQL_MAX_IDLE_CONN] != "" {
		maxIdleConns, err := typecasting.ConvertStringToInt(envCfg[ENV_POSTGRESQL_MAX_IDLE_CONN],
			ENV_POSTGRESQL_MAX_IDLE_CONN)
		if err != nil {
			return cfg, err
		} else {
			ConfigDefault.MaxIdleConns = maxIdleConns
		}
	}

	if envCfg[ENV_POSTGRESQL_CONN_MAX_IDLE_TIME] != "" {
		if maxLifetime, err := typecasting.ConvertStringToDuration(envCfg[ENV_POSTGRESQL_CONN_MAX_IDLE_TIME],
			ENV_POSTGRESQL_CONN_MAX_IDLE_TIME); err != nil {
			return cfg, err
		} else {
			ConfigDefault.ConnMaxIdleTime = maxLifetime
		}
	}

	if envCfg[ENV_POSTGRESQL_CONN_MAX_LIFE_TIME] != "" {
		if maxLifeIdleTime, err := typecasting.ConvertStringToDuration(envCfg[ENV_POSTGRESQL_CONN_MAX_LIFE_TIME],
			ENV_POSTGRESQL_CONN_MAX_LIFE_TIME); err != nil {
		} else {
			ConfigDefault.ConnMaxLifetime = maxLifeIdleTime
		}
	}

	cfg = ConfigDefault

	if len(config) > 0 {
		cfg = config[0]
	}

	return cfg, nil
}
