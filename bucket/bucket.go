package tb

import (
	"time"
)

type Bucket struct {
	LastEvent time.Time
	Tokens    []int64 // nanosecond timestamps
}
