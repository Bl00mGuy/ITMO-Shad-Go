//go:build !solution

package dupcall

import (
	"context"
	"sync"
)

type Call struct {
	mutex   sync.Mutex
	channel chan struct{}
	result  interface{}
	err     error
}

func (call *Call) startCb(context context.Context, cb func(context.Context) (interface{}, error)) {
	go func() {
		call.storeResult(context, cb)
		call.finishExecution()
	}()
}

func (call *Call) storeResult(context context.Context, cb func(context.Context) (interface{}, error)) {
	call.result, call.err = cb(context)
}

func (call *Call) finishExecution() {
	call.mutex.Lock()
	defer call.mutex.Unlock()
	close(call.channel)
	call.channel = nil
}

func (call *Call) initializeCall(context context.Context, cb func(context.Context) (interface{}, error)) chan struct{} {
	call.mutex.Lock()
	defer call.mutex.Unlock()
	if call.channel == nil {
		call.channel = make(chan struct{})
		call.startCb(context, cb)
	}
	return call.channel
}

func (call *Call) waitForResult(context context.Context, lch chan struct{}) (interface{}, error) {
	select {
	case <-lch:
		call.mutex.Lock()
		defer call.mutex.Unlock()
		return call.result, call.err
	case <-context.Done():
		return nil, context.Err()
	}
}

func (call *Call) Do(context context.Context, cb func(context.Context) (interface{}, error)) (interface{}, error) {
	lch := call.initializeCall(context, cb)
	return call.waitForResult(context, lch)
}
