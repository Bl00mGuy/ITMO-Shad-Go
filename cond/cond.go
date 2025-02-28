//go:build !solution

package cond

type Locker interface {
	Lock()
	Unlock()
}

type Cond struct {
	L     Locker
	which chan chan struct{}
}

func New(lock Locker) *Cond {
	return &Cond{
		L:     lock,
		which: make(chan chan struct{}, 200),
	}
}

func (cond *Cond) Wait() {
	defer cond.L.Lock()

	waitChannel := cond.createWaitChannel()
	cond.enqueue(waitChannel)
	cond.L.Unlock()

	cond.waitForSignal(waitChannel)
}

func (cond *Cond) createWaitChannel() chan struct{} {
	return make(chan struct{}, 1)
}

func (cond *Cond) enqueue(waitChannel chan struct{}) {
	cond.which <- waitChannel
}

func (cond *Cond) waitForSignal(waitChannel chan struct{}) {
	<-waitChannel
}

func (cond *Cond) Signal() {
	cond.signalOne()
}

func (cond *Cond) signalOne() bool {
	select {
	case waitChannel := <-cond.which:
		cond.notify(waitChannel)
		return true
	default:
		return false
	}
}

func (cond *Cond) notify(waitChannel chan struct{}) {
	select {
	case waitChannel <- struct{}{}:
	default:
		return
	}
}

func (cond *Cond) Broadcast() {
	cond.signalAll()
}

func (cond *Cond) signalAll() {
	for {
		if !cond.signalOne() {
			return
		}
	}
}
