package ds

import (
	"sync"
	"time"
)

//----- Expiry Cache Implementation -----//

const NoExpiry = -1

type CacheWithExpiry[T any] struct {
	localKeyValueCache map[string]*ExpiryCacheEntry[T]
	expiryQuanta       int64
	mutex              sync.Mutex
}

type ExpiryCacheEntry[T any] struct {
	key        string
	value      *T
	expiryTime int64
}

func GetCacheWithExpiry[T any](expiryQuanta int64) *CacheWithExpiry[T] {

	lruStore := &CacheWithExpiry[T]{
		localKeyValueCache: make(map[string]*ExpiryCacheEntry[T]),
		expiryQuanta:       expiryQuanta,
	}

	return lruStore
}

func (lruStore *CacheWithExpiry[T]) Put(key string, value *T) {
	lruStore.mutex.Lock()
	defer lruStore.mutex.Unlock()
	var expiry int64 = NoExpiry
	if lruStore.expiryQuanta != NoExpiry {
		expiry = time.Now().Unix() + lruStore.expiryQuanta
	}
	lruStore.localKeyValueCache[key] = &ExpiryCacheEntry[T]{key, value, expiry}
}

func (lruStore *CacheWithExpiry[T]) Get(key string) (*T, bool) {
	lruStore.mutex.Lock()
	defer lruStore.mutex.Unlock()
	if elem, ok := lruStore.localKeyValueCache[key]; ok {

		// return the value if it has not expired
		if elem.expiryTime == NoExpiry || elem.expiryTime < time.Now().Unix() {
			return elem.value, true
		}

		// Delete the key if it has expired
		if elem.expiryTime != NoExpiry {
			delete(lruStore.localKeyValueCache, elem.key)
		}
	}
	return nil, false
}
