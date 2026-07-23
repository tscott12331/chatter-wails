package util

import "sync"

type MutexValue[T any] struct {
	value T
	mutex  sync.Mutex
}

func NewMutexValue[T any](initial T) *MutexValue[T] {
	return &MutexValue[T]{
		value: initial,
	}
}

func (mv *MutexValue[T]) Write(value T) {
	mv.mutex.Lock()
	mv.value = value
	mv.mutex.Unlock()
}

func (mv *MutexValue[T]) Read() T {
	mv.mutex.Lock()
	value := mv.value
	mv.mutex.Unlock()

	return value
}

func (mv *MutexValue[T]) Update(updateFunc func(*T)) {
	mv.mutex.Lock()
	updateFunc(&mv.value)
	mv.mutex.Unlock()
}



type SingleWriteMutex[T any] struct {
	mutex sync.Mutex
	value T
	written bool
}

func (sfm *SingleWriteMutex[T]) IsWritten() bool {
	sfm.mutex.Lock()
	written := sfm.written
	sfm.mutex.Unlock()

	return written
}


func (sfm *SingleWriteMutex[T]) Read() (T, bool) {
	sfm.mutex.Lock()
	value := sfm.value
	written := sfm.written
	sfm.mutex.Unlock()

	return value, written
}

// returns true if first write, false if not
func (sfm *SingleWriteMutex[T]) Write(value T) bool {
	sfm.mutex.Lock()
	written := sfm.written
	if written {
		return false
	}

	sfm.value = value
	sfm.written = true
	sfm.mutex.Unlock()

	return true
}

type FetchCache[K comparable, V any] struct {
	cache MutexValue[map[K]*V]
}

func (fc *FetchCache[K, V]) Fetch(fn func() V, key K) *V {
	var result *V
	fc.cache.Update(func(m *map[K]*V) {
		value, exists := (*m)[key]
		if exists {
			result = value
			return
		}

		fetchRes := fn()

		(*m)[key] = &fetchRes
		result = &fetchRes
	})

	return result
}



type RWValue[V any] struct{
	mutex sync.RWMutex
	val V
}

func NewRWValue[V any]() *RWValue[V] {
	return &RWValue[V]{}
}

func (rwv *RWValue[V]) Set(val V) {
	rwv.mutex.Lock()
	defer rwv.mutex.Unlock()

	rwv.val = val
}

func (rwv *RWValue[V]) Get() V {
	rwv.mutex.RLock()
	defer rwv.mutex.RUnlock()

	return rwv.val
}

func (rwv *RWValue[V]) Update(fn func(*V)) {
	rwv.mutex.Lock()
	defer rwv.mutex.Unlock()

	fn(&rwv.val)
}


type RWMap[K comparable, V any] struct{
	mutex sync.RWMutex
	data map[K]V
}

func NewRWMap[K comparable, V any]() *RWMap[K,V] {
	return &RWMap[K, V]{
		data: map[K]V{},
	}
}

func (rwm *RWMap[K, V]) Get(key K) (V, bool) {
	rwm.mutex.RLock()
	defer rwm.mutex.RUnlock()

	val, exists := rwm.data[key]
	return val, exists
}

func (rwm *RWMap[K, V]) Set(key K, val V) {
	rwm.mutex.Lock()
	defer rwm.mutex.Unlock()

	rwm.data[key] = val
}

func (rwm *RWMap[K, V]) Delete(key K) {
	rwm.mutex.Lock()
	defer rwm.mutex.Unlock()

	delete(rwm.data, key)
}
