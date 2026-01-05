package store

import "sync"

type Repository[K comparable, V any] struct {
	mu   sync.RWMutex
	data map[K]V
}

// creates  new generic repository
func NewRepository[K comparable, V any]() *Repository[K, V] {
	return &Repository[K, V]{
		data: make(map[K]V),
	}
}

// set value with the given key
func (r *Repository[K, V]) Set(key K, value V) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.data[key] = value
}

// get retrieves  value from key
func (r *Repository[K, V]) Get(key K) (V, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	val, exists := r.data[key]
	return val, exists
}

// ret all values
func (r *Repository[K, V]) GetAll() []V {
	r.mu.RLock()
	defer r.mu.RUnlock()

	values := make([]V, 0, len(r.data))
	for _, v := range r.data {
		values = append(values, v)
	}
	return values
}

func (r *Repository[K, V]) Update(key K, updateFn func(V) V) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	if val, exists := r.data[key]; exists {
		r.data[key] = updateFn(val)
		return true
	}
	return false
}

// Count() ret number of items
func (r *Repository[K, V]) Count() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.data)
}
