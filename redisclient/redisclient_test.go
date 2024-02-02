package redisclient_test

import (
	"context"
	"testing"

	env "github.com/joho/godotenv"
	"github.com/rohanraj7316/utils/redisclient"
)

const (
	envFile = ".env"

	key   = "testKey"
	value = "testValue"
)

var loadEnv = env.Load

func TestPing(t *testing.T) {
	ctx := context.Background()

	// load env config
	err := loadEnv(envFile)
	if err != nil {
		t.Error(err)
	}

	cache, err := redisclient.New()
	if err != nil {
		t.Errorf("failed to load cache: %s", err)
	}

	_, err = cache.Client.Ping(ctx).Result()
	if err != nil {
		t.Errorf("failed to ping redis: %s", err)
	}
}

func TestGet(t *testing.T) {
	ctx := context.Background()

	// load env config
	err := loadEnv(envFile)
	if err != nil {
		t.Error(err)
	}

	cache, err := redisclient.New()
	if err != nil {
		t.Errorf("failed to load cache: %s", err)
	}

	_, err = cache.Get(ctx, key).Result()
	if err != nil {
		t.Errorf("failed to fetch data: %s", err)
	}
}

func TestSet(t *testing.T) {
	ctx := context.Background()

	// load env config
	err := loadEnv(envFile)
	if err != nil {
		t.Error(err)
	}

	cache, err := redisclient.New()
	if err != nil {
		t.Errorf("failed to load cache: %s", err)
	}

	if err := cache.Set(ctx, "testKey", "testValue", 0).Err(); err != nil {
		t.Errorf("failed to set cache: %s", err)
	}
}

func TestDel(t *testing.T) {

}
