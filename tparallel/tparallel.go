//go:build !solution

package tparallel

type T struct {
	finished chan struct{}
	barrier  chan struct{}
	parallel bool
	parent   *T
	children []*T
}

func newT(parent *T) *T {
	return &T{
		finished: make(chan struct{}),
		barrier:  make(chan struct{}),
		parallel: false,
		parent:   parent,
		children: make([]*T, 0),
	}
}

func (t *T) Parallel() {
	t.parallel = true
	t.parent.children = append(t.parent.children, t)

	t.finished <- struct{}{}
	<-t.parent.barrier
}

func (t *T) waitForChildren() {
	if len(t.children) > 0 {
		close(t.barrier)

		for _, child := range t.children {
			<-child.finished
		}
	}
}

func (t *T) signalParentCompletion() {
	if t.parallel {
		t.parent.finished <- struct{}{}
	}
}

func (t *T) runner(subtest func(t *T)) {
	subtest(t)
	t.waitForChildren()
	t.signalParentCompletion()

	t.finished <- struct{}{}
}

func (t *T) Run(subtest func(t *T)) {
	child := newT(t)
	go child.runner(subtest)
	<-child.finished
}

func Run(topTests []func(t *T)) {
	root := newT(nil)
	for _, subtest := range topTests {
		root.Run(subtest)
	}

	close(root.barrier)

	if len(root.children) > 0 {
		<-root.finished
	}
}
