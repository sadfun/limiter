package storage

import (
	tb "github.com/sadfun/limiter/bucket"
	"sync"
)

type MapStorage[K comparable] struct {
	buckets map[K]tb.Bucket
	locks   map[K]*sync.Mutex
	mx      sync.RWMutex
}

func NewMapStorage[K comparable]() Storage[K] {
	return &MapStorage[K]{
		buckets: make(map[K]tb.Bucket),
		locks:   make(map[K]*sync.Mutex),
	}
}

func (ms *MapStorage[K]) getLock(k K) *sync.Mutex {
	ms.mx.Lock()
	defer ms.mx.Unlock()

	if _, ok := ms.locks[k]; !ok {
		ms.locks[k] = &sync.Mutex{}
	}

	return ms.locks[k]
}

func (ms *MapStorage[K]) Update(key K, updater func(bucket tb.Bucket) (newBucket tb.Bucket)) {
	l := ms.getLock(key)

	// We must unlock the lock not only after the updater is done, but also around read-write from the map
	l.Lock()
	defer l.Unlock() // User may panic in updater

	ms.mx.RLock()
	bucket := ms.buckets[key]
	ms.mx.RUnlock()

	bucket = updater(bucket)

	ms.mx.Lock()
	ms.buckets[key] = bucket
	ms.mx.Unlock()
}
