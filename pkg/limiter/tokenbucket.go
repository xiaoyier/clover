package limiter

import (
	"sync"
	"time"
)

var (
	bucket *Bucket
	once   sync.Once
)

func Init(capacity int64) *Bucket {

	once.Do(func() {
		bucket = NewBucketWithQuantum(time.Millisecond, capacity, 1)
	})

	return bucket
}

func GetBucket() *Bucket {
	return bucket
}
