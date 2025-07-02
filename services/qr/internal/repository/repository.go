package repository

import (
	"context"
	"time"
)

type QrRepository interface {
	Delete(ctx context.Context, key string) error
	Set(ctx context.Context, key string, value any, ttl time.Duration) error
	Get(ctx context.Context, key string) (string, error)
}
