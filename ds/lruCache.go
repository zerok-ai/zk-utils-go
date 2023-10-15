package ds

import (
	"container/list"
)

type Cache[T any] interface {
	Put(key string, value *T)
	Get(key string) (*T, bool)
}

//----- LRU Cache Implementation -----//

type LRUCache[T any] struct {
	localKeyValueCache map[string]*list.Element
	capacity           int
	recencyList        *list.List
}

type Entry[T any] struct {
	key   string
	value *T
}

func GetLRUCache[T any](cacheSize int) *LRUCache[T] {

	lruStore := &LRUCache[T]{
		localKeyValueCache: make(map[string]*list.Element),
		capacity:           cacheSize,
		recencyList:        list.New(),
	}

	return lruStore
}

func (lruStore *LRUCache[T]) Put(key string, value *T) {
	if elem, ok := lruStore.localKeyValueCache[key]; ok {
		elem.Value.(*Entry[T]).value = value
		lruStore.recencyList.MoveToFront(elem)
		return
	}

	if len(lruStore.localKeyValueCache) >= lruStore.capacity {
		oldest := lruStore.recencyList.Back()
		if oldest != nil {
			delete(lruStore.localKeyValueCache, oldest.Value.(*Entry[T]).key)
			lruStore.recencyList.Remove(oldest)
		}
	}

	newEntry := &Entry[T]{key, value}
	newElem := lruStore.recencyList.PushFront(newEntry)
	lruStore.localKeyValueCache[key] = newElem
}

func (lruStore *LRUCache[T]) Get(key string) (*T, bool) {
	if elem, ok := lruStore.localKeyValueCache[key]; ok {
		lruStore.recencyList.MoveToFront(elem)
		return elem.Value.(*Entry[T]).value, true
	}
	return nil, false
}
