//go:build !solution

package waitgroup

type WaitGroup struct {
	run     chan struct{}
	counter chan int
}

func New() *WaitGroup {
	wg := &WaitGroup{
		run:     nil,
		counter: make(chan int, 1),
	}
	wg.counter <- 0
	return wg
}

func (wg *WaitGroup) Add(delta int) {
	counter := <-wg.counter

	if counter == 0 && delta > 0 {
		wg.run = make(chan struct{})
	}

	newCounter := counter + delta
	if newCounter < 0 {
		panic("negative WaitGroup counter")
	}

	wg.counter <- newCounter

	if newCounter == 0 {
		close(wg.run)
	}
}

func (wg *WaitGroup) Done() {
	wg.Add(-1)
}

func (wg *WaitGroup) Wait() {
	if wg.run == nil {
		return
	}
	<-wg.run
}
