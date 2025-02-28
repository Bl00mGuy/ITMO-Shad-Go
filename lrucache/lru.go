//go:build !solution

package lrucache

type LRU struct {
	cache    map[int]*data
	head     *data
	tail     *data
	capacity int
	size     int
}

type data struct {
	key   int
	value int
	prev  *data
	next  *data
}

func New(cap int) *LRU {
	return &LRU{
		cache:    make(map[int]*data, cap),
		capacity: cap,
	}
}

func (cache *LRU) Get(key int) (int, bool) {
	if elem, found := cache.cache[key]; found {
		cache.moveToBack(elem)
		return elem.value, true
	}
	return 0, false
}

func (cache *LRU) Set(key, value int) {
	if elem, found := cache.cache[key]; found {
		elem.value = value
		cache.moveToBack(elem)
		return
	}
	newElem := &data{key: key, value: value}
	cache.cache[key] = newElem
	cache.addToBack(newElem)
	cache.size++
	if cache.size > cache.capacity {
		removed := cache.removeFront()
		delete(cache.cache, removed.key)
		cache.size--
	}
}

func (cache *LRU) Range(f func(key, value int) bool) {
	for elem := cache.head; elem != nil; elem = elem.next {
		if !f(elem.key, elem.value) {
			break
		}
	}
}

func (cache *LRU) Clear() {
	cache.cache = make(map[int]*data, cache.capacity)
	cache.head = nil
	cache.tail = nil
	cache.size = 0
}

func (cache *LRU) moveToBack(elem *data) {
	if cache.tail == elem {
		return
	}
	cache.remove(elem)
	cache.addToBack(elem)
}

func (cache *LRU) addToBack(elem *data) {
	if cache.tail != nil {
		cache.tail.next = elem
		elem.prev = cache.tail
	} else {
		cache.head = elem
	}
	cache.tail = elem
}

func (cache *LRU) removeFront() *data {
	if cache.head == nil {
		return nil
	}
	removed := cache.head
	cache.remove(removed)
	return removed
}

func (cache *LRU) remove(elem *data) {
	if elem.prev != nil {
		elem.prev.next = elem.next
	} else {
		cache.head = elem.next
	}
	if elem.next != nil {
		elem.next.prev = elem.prev
	} else {
		cache.tail = elem.prev
	}
	elem.prev = nil
	elem.next = nil
}
