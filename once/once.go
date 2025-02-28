//go:build !solution

package once

type Once struct {
	runner chan struct{}
	done   chan struct{}
}

func New() *Once {
	return &Once{
		runner: make(chan struct{}, 1),
		done:   make(chan struct{}),
	}
}

func (once *Once) Do(f func()) {
	if once.start() {
		defer once.finish()
		f()
	} else {
		once.wait()
	}
}

func (once *Once) start() bool {
	select {
	case once.runner <- struct{}{}:
		return true
	default:
		return false
	}
}

func (once *Once) finish() {
	close(once.done)
}

func (once *Once) wait() {
	<-once.done
}
