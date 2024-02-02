package mongodb

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"strconv"
	"time"

	"github.com/pkg/errors"

	"github.com/rohanraj7316/logger"
	"github.com/rohanraj7316/utils/env"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	ENV_MONGO_DB_CONNECTION_STRING = "MONGO_DB_CONNECTION_STRING"
	ENV_MONGO_DB_DATABASE_NAME     = "MONGO_DB_DATABASE_NAME"
	ENV_MONGO_DB_MAX_IDLE_TIME     = "MONGO_DB_MAX_IDLE_TIME"
	ENV_MONGO_DB_MAX_POOL_SIZE     = "MONGO_DB_MAX_POOL_SIZE"
	ENV_MONGO_DB_MIN_POOL_SIZE     = "MONGO_DB_MIN_POOL_SIZE"
	ENV_MONGO_DB_IS_TLS            = "MONGO_DB_IS_TLS"
	ENV_MONGO_DB_TLS_PEM_FILE_PATH = "MONGO_DB_TLS_PEM_FILE_PATH"
	ENV_MONGO_DB_TRACE_ENABLE      = "MONGO_DB_TRACE_ENABLE"
)

type Config struct {
	cOptions    *options.ClientOptions
	dOptions    *options.DatabaseOptions
	collOptions *options.CollectionOptions

	DbName       string
	Client       *mongo.Client
	IsValidation bool
}

var ConfigDefault = Config{
	cOptions:     &options.ClientOptions{},
	collOptions:  &options.CollectionOptions{},
	IsValidation: true,
}

func getCustomTLSConfig(caFile string) (*tls.Config, error) {
	tlsConfig := new(tls.Config)
	certs, err := ioutil.ReadFile(caFile)
	if err != nil {
		return tlsConfig, errors.WithStack(err)
	}

	tlsConfig.RootCAs = x509.NewCertPool()
	ok := tlsConfig.RootCAs.AppendCertsFromPEM(certs)
	if !ok {
		return tlsConfig, errors.Errorf("failed parsing pem file")
	}

	return tlsConfig, nil
}

func configDefault(config ...Config) (Config, error) {
	if len(config) < 1 {
		rEnvCfg := []string{
			ENV_MONGO_DB_CONNECTION_STRING,
			ENV_MONGO_DB_DATABASE_NAME,
		}
		uEnvCfg := []string{
			ENV_MONGO_DB_MAX_IDLE_TIME,
			ENV_MONGO_DB_MIN_POOL_SIZE,
			ENV_MONGO_DB_MAX_POOL_SIZE,
			ENV_MONGO_DB_IS_TLS,
			ENV_MONGO_DB_TLS_PEM_FILE_PATH,
			ENV_MONGO_DB_TRACE_ENABLE,
		}
		envCfg := env.EnvData(rEnvCfg, uEnvCfg)

		ConfigDefault.cOptions.
			ApplyURI(envCfg[ENV_MONGO_DB_CONNECTION_STRING]).
			SetBSONOptions(&options.BSONOptions{
				UseJSONStructTags: true,
				OmitZeroStruct:    true,
			})

		if _, ok := envCfg[ENV_MONGO_DB_MAX_IDLE_TIME]; ok {
			maxConnIdleTime, err := time.ParseDuration(fmt.Sprintf("%ss", envCfg[ENV_MONGO_DB_MAX_IDLE_TIME]))
			if err != nil {
				return ConfigDefault, errors.WithStack(err)
			}
			ConfigDefault.cOptions = ConfigDefault.cOptions.SetMaxConnIdleTime(maxConnIdleTime)
		}

		if _, ok := envCfg[ENV_MONGO_DB_MAX_POOL_SIZE]; ok {
			maxPoolSize, err := strconv.Atoi(envCfg[ENV_MONGO_DB_MAX_POOL_SIZE])
			if err != nil {
				return ConfigDefault, errors.WithStack(err)
			}
			ConfigDefault.cOptions = ConfigDefault.cOptions.SetMaxPoolSize(uint64(maxPoolSize))
		}

		if _, ok := envCfg[ENV_MONGO_DB_MIN_POOL_SIZE]; ok {
			minPoolSize, err := strconv.Atoi(envCfg[ENV_MONGO_DB_MIN_POOL_SIZE])
			if err != nil {
				return ConfigDefault, errors.WithStack(err)
			}
			ConfigDefault.cOptions = ConfigDefault.cOptions.SetMinPoolSize(uint64(minPoolSize))
		}

		if isTls, ok := envCfg[ENV_MONGO_DB_IS_TLS]; ok {
			if isTls == "true" {
				if filepath, ok := envCfg[ENV_MONGO_DB_TLS_PEM_FILE_PATH]; ok {
					tlsCfg, err := getCustomTLSConfig(filepath)
					if err != nil {
						return ConfigDefault, errors.WithStack(err)
					}

					ConfigDefault.cOptions = ConfigDefault.cOptions.SetTLSConfig(tlsCfg)
				} else {
					return ConfigDefault, errors.Errorf("empty tls pathname")
				}
			}
		}

		if _, ok := envCfg[ENV_MONGO_DB_DATABASE_NAME]; ok {
			ConfigDefault.DbName = envCfg[ENV_MONGO_DB_DATABASE_NAME]
		}

		if isTraceEnable, ok := envCfg[ENV_MONGO_DB_TRACE_ENABLE]; ok {
			if isTraceEnable == "true" {
				ConfigDefault.cOptions.SetMonitor(&event.CommandMonitor{
					Started: func(ctx context.Context, cse *event.CommandStartedEvent) {
						fields := []logger.Field{
							{
								Key:   "rawQuery",
								Value: cse.Command.String(),
							},
							// sync our request id with this trace id
							{
								Key:   "requestId",
								Value: cse.RequestID,
							},
							{
								Key:   "database",
								Value: cse.DatabaseName,
							},
							{
								Key:   "serviceId",
								Value: cse.ServiceID,
							},
						}

						logger.Info("mongodb trace log", fields...)
					},
					Succeeded: func(ctx context.Context, cse *event.CommandSucceededEvent) {
						// NOTE: implement in future
					},
					Failed: func(ctx context.Context, cfe *event.CommandFailedEvent) {
						fields := []logger.Field{
							{
								Key:   "requestId",
								Value: cfe.CommandFinishedEvent.RequestID,
							},
							{
								Key:   "latency",
								Value: cfe.CommandFinishedEvent.DurationNanos / int64(time.Millisecond),
							},
						}

						logger.Error(cfe.Failure, fields...)
					},
				})
			}
		}
	}

	return ConfigDefault, nil
}
