package redislimiter

import (
	"context"
	"log"
	"time"

	"github.com/go-redis/redis_rate/v10"
	"github.com/redis/go-redis/v9"
)

type LimiterInterface interface {
	Allow(ctx context.Context, key string, rate int, period time.Duration) (bool, error)
	Reset(ctx context.Context, key string) error
}
type Limiter struct {
	l *redis_rate.Limiter
}

var _ LimiterInterface = (*Limiter)(nil)

func NewLimiter(client *redis.Client) *Limiter {
	return &Limiter{l: redis_rate.NewLimiter(client)}
}

func (r *Limiter) Allow(ctx context.Context, key string, rate int, period time.Duration) (bool, error) {
	op := "redisLimiter.Allow"
	res, err := r.l.Allow(ctx, key, redis_rate.Limit{
		Rate:   rate,
		Burst:  rate,
		Period: period,
	})
	if err != nil {
		log.Println("op", op, "Error:", err)
		return false, err
	}
	return res.Allowed > 0, nil
}

func (r *Limiter) Reset(ctx context.Context, key string) error {
	op := "redisLimiter.reset"
	err := r.l.Reset(ctx, key)
	if err != nil {
		log.Println("op:", op, "rate limiter error:", err)
		return err
	}
	return nil
}
