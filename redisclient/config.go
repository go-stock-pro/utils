package redisclient

import (
	"context"
	"fmt"
	"strconv"

	"github.com/go-redis/redis/v8"
	"github.com/rohanraj7316/utils/env"
)

const (
	ENV_REDIS_RUNTIME_ENV = "REDIS_ENV"
	ENV_REDIS_HOST        = "REDIS_HOST"
	ENV_REDIS_PORT        = "REDIS_PORT"
	ENV_REDIS_AUTH        = "REDIS_AUTH"
	ENV_REDIS_PREFIX      = "REDIS_PREFIX"
	ENV_REDIS_POOL_SIZE   = "REDIS_POOL_SIZE"
	ENV_REDIS_MAX_RETRIES = "REDIS_MAX_RETRIES"
)

type Config struct {
	Redis  *redis.Options
	Prefix string
}

var ConfigDefault = Config{
	Redis: &redis.Options{
		OnConnect: func(ctx context.Context, cn *redis.Conn) error {
			return nil
		},
	},
}

func configDefault(config ...Config) (Config, error) {
	if len(config) < 1 {
		// extracting configs from env
		rConfigs := []string{
			ENV_REDIS_HOST,
			ENV_REDIS_PORT,
			ENV_REDIS_PREFIX,
		}
		uConfigs := []string{
			ENV_REDIS_AUTH,
			ENV_REDIS_POOL_SIZE,
			ENV_REDIS_MAX_RETRIES,
			ENV_REDIS_RUNTIME_ENV,
		}
		eConfig := env.EnvData(rConfigs, uConfigs)

		ConfigDefault.Redis.Addr = fmt.Sprintf("%s:%s", eConfig[ENV_REDIS_HOST], eConfig[ENV_REDIS_PORT])

		if eAuth := eConfig[ENV_REDIS_AUTH]; len(eAuth) != 0 {
			ConfigDefault.Redis.Password = eConfig[ENV_REDIS_AUTH]
			// if eConfig[ENV_REDIS_RUNTIME_ENV] != "local" {
			// 	ConfigDefault.Redis.TLSConfig = &tls.Config{}
			// }
		}

		// max retries
		if eMaxTries := eConfig[ENV_REDIS_MAX_RETRIES]; len(eMaxTries) != 0 {
			mRetries, err := strconv.Atoi(eMaxTries)
			if err != nil {
				return Config{}, err
			}
			ConfigDefault.Redis.MaxRetries = mRetries
		}

		// pool size
		if ePoolSize := eConfig[ENV_REDIS_POOL_SIZE]; len(ePoolSize) != 0 {
			pSize, err := strconv.Atoi(eConfig[ENV_REDIS_POOL_SIZE])
			if err != nil {
				return Config{}, err
			}
			ConfigDefault.Redis.PoolSize = pSize
		}

		ConfigDefault.Prefix = eConfig[ENV_REDIS_PREFIX]

		return ConfigDefault, nil
	}

	cfg := config[0]

	// cfg.Redis.TLSConfig = &tls.Config{}

	return cfg, nil
}
