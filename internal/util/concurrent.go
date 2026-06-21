package util

import "sync"

type MutexValue[T any] struct {
	value T
	mutex  sync.Mutex
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
