//go:build !solution

package consistenthash

import (
	"crypto/sha256"
	"fmt"
	"sort"
	"sync"
)

type Node interface {
	ID() string
}

type ConsistentHash[N Node] struct {
	virtualReplicas int
	hashRing        []int
	hashToNode      map[int]*N
	mutex           sync.RWMutex
}

func New[N Node]() *ConsistentHash[N] {
	return &ConsistentHash[N]{
		hashRing:        []int{},
		virtualReplicas: 92,
		hashToNode:      make(map[int]*N),
	}
}

func (h *ConsistentHash[N]) AddNode(node *N) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	for i := 0; i < h.virtualReplicas; i++ {
		virtualNodeID := generateVirtualNodeID((*node).ID(), i)
		hash := hashKey(virtualNodeID)
		h.addHashToRing(hash, node)
	}
}

func (h *ConsistentHash[N]) GetNode(key string) *N {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	if len(h.hashRing) == 0 {
		return nil
	}

	hash := hashKey(key)
	index := findClosestHash(h.hashRing, hash)
	return h.hashToNode[h.hashRing[index]]
}

func (h *ConsistentHash[N]) RemoveNode(node *N) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	for i := 0; i < h.virtualReplicas; i++ {
		virtualNodeID := generateVirtualNodeID((*node).ID(), i)
		hash := hashKey(virtualNodeID)
		h.removeHashFromRing(hash)
	}
}

func generateVirtualNodeID(nodeID string, replica int) string {
	return fmt.Sprintf("%s#%d", nodeID, replica)
}

func hashKey(key string) int {
	hash := sha256.Sum256([]byte(key))
	return int(hash[0])<<24 | int(hash[1])<<16 | int(hash[2])<<8 | int(hash[3])
}

func (h *ConsistentHash[N]) addHashToRing(hash int, node *N) {
	h.hashRing = insertInOrder(h.hashRing, hash)
	h.hashToNode[hash] = node
}

func (h *ConsistentHash[N]) removeHashFromRing(hash int) {
	h.hashRing = removeFromOrder(h.hashRing, hash)
	delete(h.hashToNode, hash)
}

func findClosestHash(ring []int, hash int) int {
	index := sort.Search(len(ring), func(i int) bool { return ring[i] >= hash })
	if index == len(ring) {
		return 0
	}
	return index
}

func insertInOrder(slice []int, value int) []int {
	index := sort.Search(len(slice), func(i int) bool { return slice[i] >= value })
	slice = append(slice, 0)
	copy(slice[index+1:], slice[index:])
	slice[index] = value
	return slice
}

func removeFromOrder(slice []int, value int) []int {
	index := sort.Search(len(slice), func(i int) bool { return slice[i] >= value })
	if index < len(slice) && slice[index] == value {
		return append(slice[:index], slice[index+1:]...)
	}
	return slice
}
