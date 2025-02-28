//go:build !solution

package rwmutex

type RWMutex struct {
	read  chan int
	write chan struct{}
}

func New() *RWMutex {
	rw := &RWMutex{
		read:  make(chan int, 1),
		write: make(chan struct{}, 1),
	}
	rw.read <- 0
	rw.write <- struct{}{}
	return rw
}

func (rw *RWMutex) RLock() {
	rw.acquireReadLock()
}

func (rw *RWMutex) RUnlock() {
	rw.releaseReadLock()
}

func (rw *RWMutex) Lock() {
	rw.acquireWriteLock()
}

func (rw *RWMutex) Unlock() {
	rw.releaseWriteLock()
}

func (rw *RWMutex) acquireReadLock() {
	counter := <-rw.read
	if counter == 0 {
		<-rw.write
	}
	rw.read <- counter + 1
}

func (rw *RWMutex) releaseReadLock() {
	counter := <-rw.read
	if counter == 1 {
		rw.write <- struct{}{}
	}
	rw.read <- counter - 1
}

func (rw *RWMutex) acquireWriteLock() {
	<-rw.write
}

func (rw *RWMutex) releaseWriteLock() {
	rw.write <- struct{}{}
}
