# Redisclient

**redisclient** built to establish connection with redis.

## Signature

```
func New(config ...Config) (Config, error)
```

## Config

```
type Config struct {
	Redis  *redis.Options
	Client *redis.Client
	Prefix string
}
```
## Env-variable

```
export REDIS_HOST
export REDIS_PORT
export REDIS_AUTH
export REDIS_PREFIX
export REDIS_MAX_RETRIES
export REDIS_POOL_SIZE
```