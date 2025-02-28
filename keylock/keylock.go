//go:build !solution

package keylock

import (
	"sort"
	"sync"
)

type signalChan = chan struct{}

type KeyLock struct {
	mutex sync.Mutex
	locks map[string]signalChan
}

type preparedKeys struct {
	sortedKeys []string
	cancelled  bool
}

func New() *KeyLock {
	return &KeyLock{
		locks: make(map[string]signalChan),
	}
}

func prepareAndSortKeys(keys []string) preparedKeys {
	sortedKeys := make([]string, len(keys))
	copy(sortedKeys, keys)
	sort.Strings(sortedKeys)

	return preparedKeys{
		sortedKeys: sortedKeys,
		cancelled:  false,
	}
}

func (keyLock *KeyLock) acquireAllKeys(prepKeys preparedKeys, cancel <-chan struct{}) bool {
	for _, key := range prepKeys.sortedKeys {
		lockChan := keyLock.obtainLockChannel(key)
		if !keyLock.waitForKeyRelease(lockChan, cancel) {
			prepKeys.cancelled = true
			break
		}
	}

	return prepKeys.cancelled
}

func (keyLock *KeyLock) obtainLockChannel(key string) signalChan {
	keyLock.mutex.Lock()
	defer keyLock.mutex.Unlock()

	if ch, exists := keyLock.locks[key]; exists {
		return ch
	}

	ch := make(signalChan, 1)
	ch <- struct{}{}
	keyLock.locks[key] = ch

	return ch
}

func (keyLock *KeyLock) waitForKeyRelease(lockChan signalChan, cancel <-chan struct{}) bool {
	select {
	case <-cancel:
		return false
	case <-lockChan:
		return true
	}
}

func (keyLock *KeyLock) releaseKeys(keys []string) {
	keyLock.mutex.Lock()
	defer keyLock.mutex.Unlock()

	for _, key := range keys {
		if lockChan, exists := keyLock.locks[key]; exists {
			select {
			case lockChan <- struct{}{}:
			default:
			}
		}
	}
}

func (keyLock *KeyLock) LockKeys(keys []string, cancel <-chan struct{}) (bool, func()) {
	prepKeys := prepareAndSortKeys(keys)

	if keyLock.acquireAllKeys(prepKeys, cancel) {
		keyLock.releaseKeys(prepKeys.sortedKeys)
		return true, func() {}
	}

	return false, func() { keyLock.releaseKeys(prepKeys.sortedKeys) }
}
