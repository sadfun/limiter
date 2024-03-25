package limiter

import (
	tb "github.com/sadfun/limiter/bucket"
	"github.com/sadfun/limiter/storage"
	"slices"
	"time"
)

type Limiter[K comparable] struct {
	storage storage.Storage[K]

	limit int
	t     time.Duration
}

type Config[K comparable] struct {
	Limit    int
	Duration time.Duration
	Storage  storage.Storage[K]
}

func fillDefaultConfig[K comparable](config *Config[K]) {
	if config.Storage == nil {
		config.Storage = storage.NewMapStorage[K]()
	}
}

func NewLimiter[K comparable](config *Config[K]) *Limiter[K] {
	fillDefaultConfig(config)

	return &Limiter[K]{
		storage: config.Storage,
		limit:   config.Limit,
		t:       config.Duration,
	}
}

func (limiter *Limiter[K]) Use(key K) (ok bool) {
	return limiter.UseN(key, 1)
}

func (limiter *Limiter[K]) UseN(key K, n int) (ok bool) {
	if limiter.limit < n {
		return false
	}

	now := time.Now()

	limiter.storage.Update(key, func(bucket tb.Bucket) (newBucket tb.Bucket) {
		limiter.dropExpiredTokens(now, &bucket)

		if (len(bucket.Tokens) + n) > limiter.limit {
			return bucket
		}

		slices.Grow(bucket.Tokens, limiter.limit-len(bucket.Tokens))
		for i := 0; i < n; i++ {
			bucket.Tokens = append(bucket.Tokens, now.UnixNano())
		}

		ok = true

		return bucket
	})

	return ok
}

func (limiter *Limiter[K]) dropExpiredTokens(now time.Time, tb *tb.Bucket) {
	i, _ := slices.BinarySearch(tb.Tokens, now.Add(-limiter.t).UnixNano())
	tb.Tokens = tb.Tokens[i:]
}
