package storage

import (
	tb "github.com/sadfun/limiter/bucket"
)

type Marshaler interface {
	Marshal() ([]byte, error)
}

type Storage[K comparable] interface {
	// Update the bucket with new values.
	// Must hold a distributed lock on the key.
	// If the bucket does not exist, it should be created.
	Update(key K, updater func(bucket tb.Bucket) (newBucket tb.Bucket))
}
