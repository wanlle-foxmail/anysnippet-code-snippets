package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

const releaseLockScript = `
if redis.call("GET", KEYS[1]) == ARGV[1] then
  return redis.call("DEL", KEYS[1])
end
return 0
`

type RedisLockStore interface {
	SetNX(ctx context.Context, key string, value string, ttl time.Duration) (bool, error)
	CompareAndDelete(ctx context.Context, key string, expectedValue string) (bool, error)
}

type GoRedisLockStore struct {
	client *redis.Client
}

func (store GoRedisLockStore) SetNX(ctx context.Context, key string, value string, ttl time.Duration) (bool, error) {
	return store.client.SetNX(ctx, key, value, ttl).Result()
}

func (store GoRedisLockStore) CompareAndDelete(ctx context.Context, key string, expectedValue string) (bool, error) {
	deletedCount, err := store.client.Eval(ctx, releaseLockScript, []string{key}, expectedValue).Int()
	if err != nil {
		return false, err
	}
	return deletedCount == 1, nil
}

func AcquireRedisLock(ctx context.Context, store RedisLockStore, lockKey string, ownerToken string, ttl time.Duration) (bool, error) {
	if err := validateLockInput(store, lockKey, ownerToken); err != nil {
		return false, err
	}
	if ttl <= 0 {
		return false, errors.New("ttl must be greater than 0")
	}

	// Flow:
	//   validate key, owner token, and ttl
	//      |
	//      +-> SET NX succeeds -> lock acquired
	//      `-> SET NX fails or errors -> return false or wrapped error
	acquired, err := store.SetNX(ctx, lockKey, ownerToken, ttl)
	if err != nil {
		return false, fmt.Errorf("acquire redis lock %s: %w", lockKey, err)
	}
	return acquired, nil
}

func ReleaseRedisLock(ctx context.Context, store RedisLockStore, lockKey string, ownerToken string) (bool, error) {
	if err := validateLockInput(store, lockKey, ownerToken); err != nil {
		return false, err
	}

	released, err := store.CompareAndDelete(ctx, lockKey, ownerToken)
	if err != nil {
		return false, fmt.Errorf("release redis lock %s: %w", lockKey, err)
	}
	return released, nil
}

func validateLockInput(store RedisLockStore, lockKey string, ownerToken string) error {
	if store == nil {
		return errors.New("store is required")
	}
	if strings.TrimSpace(lockKey) == "" {
		return errors.New("lock key is required")
	}
	if strings.TrimSpace(ownerToken) == "" {
		return errors.New("owner token is required")
	}
	return nil
}

func main() {
	ctx := context.Background()
	client := redis.NewClient(&redis.Options{Addr: "localhost:6379"})
	store := GoRedisLockStore{client: client}

	locked, err := AcquireRedisLock(ctx, store, "jobs:123", "worker-1", 15*time.Second)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("lock acquired=%t", locked)
}