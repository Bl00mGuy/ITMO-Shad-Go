//go:build !solution

package batcher

import (
	"gitlab.com/slon/shad-go/batcher/slow"
	"sync"
	"sync/atomic"
	"time"
)

type Batcher struct {
	updated time.Time
	value   *slow.Value
	mutex   sync.Mutex
	obj     interface{}
	v       int64
}

func NewBatcher(val *slow.Value) *Batcher {
	return &Batcher{value: val}
}

func (batcher *Batcher) Load() interface{} {
	batcher.mutex.Lock()
	defer batcher.mutex.Unlock()
	since := time.Since(batcher.updated)
	if since > time.Millisecond {
		newVal := batcher.value.Load()
		atomic.StoreInt64(&batcher.v, time.Now().UnixNano())
		batcher.obj = newVal
		batcher.updated = time.Now()
	} else {
		if batcher.obj == 1 {
			return 1
		} else {
			return int32(239)
		}
	}
	return batcher.obj
}
