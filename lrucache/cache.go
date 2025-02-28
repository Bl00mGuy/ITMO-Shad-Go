//go:build !change

package lrucache

type Cache interface {
	Get(key int) (int, bool)
	Set(key, value int)
	Range(f func(key, value int) bool)
	Clear()
}
