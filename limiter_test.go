package limiter

import (
	"github.com/sadfun/limiter/storage"
	"testing"
	"time"
)

func TestLimiter(t *testing.T) {
	t.Run("Dummy", func(t *testing.T) {
		dummyTest(t, 10, 100*time.Millisecond)
		dummyTest(t, 1, 100*time.Millisecond)
	})

	t.Run("Leak", func(t *testing.T) {
		// Fill full bucket quickly, then try to do one more request after duration/limit.
		// This test is failing with many Token Bucket implementations.
		limit := 10
		duration := time.Second

		limiter := NewLimiter[string](&Config[string]{
			Limit:    limit,
			Duration: duration,
			Storage:  storage.NewMapStorage[string](),
		})

		if !limiter.UseN("key", limit) {
			t.Error("UseN failed")
		}

		time.Sleep(
			time.Duration(duration.Nanoseconds()/int64(limit)) + time.Millisecond,
		)

		if limiter.UseN("key", 1) {
			t.Error("Limiter leaked request after divided duration")
		}

		time.Sleep(
			time.Duration(duration.Nanoseconds()-(duration.Nanoseconds()/int64(limit))) - 10*time.Millisecond,
		)

		if limiter.UseN("key", 1) {
			t.Error("Limiter leaked request after almost full duration")
		}

		time.Sleep(
			10 * time.Millisecond,
		)

		if !limiter.UseN("key", limit) {
			t.Error("Limiter did not allow requests after full duration")
		}
	})
}

func dummyTest(t *testing.T, limit int, duration time.Duration) {
	// Fill full bucket quickly, then try to do one more request without and with delay.

	limiter := NewLimiter[string](&Config[string]{
		Limit:    limit,
		Duration: duration,
		Storage:  storage.NewMapStorage[string](),
	})

	if !limiter.UseN("key", limit) {
		t.Error("UseN failed")
	}

	if limiter.UseN("key", 1) {
		t.Error("UseN succeeded")
	}

	time.Sleep(duration)

	if !limiter.UseN("key", 1) {
		t.Error("UseN failed")
	}

	if !limiter.UseN("key", limit-1) {
		t.Error("UseN failed")
	}
}
