package common

import (
	"sync"
)

type MapKeyValue[K comparable, V any] struct {
	data  map[K]V
	mutex sync.RWMutex
}

func NewMapKeyValue[K comparable, V any]() *MapKeyValue[K, V] {
	return &MapKeyValue[K, V]{
		data: make(map[K]V),
	}
}

func (m *MapKeyValue[K, V]) Get(key K) V {
	var zero V
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	v, ok := m.data[key]
	if !ok {
		return zero
	}
	return v
}

func (m *MapKeyValue[K, V]) Set(key K, value V) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.data[key] = value
}
