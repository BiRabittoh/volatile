package volatile

import (
	"fmt"
	"time"
)

type Element[V any] struct {
	value     *V
	Timestamp time.Time
}

type Volatile[K comparable, V any] struct {
	data       map[K]Element[V]
	timeToLive time.Duration
}

func (v *Volatile[K, V]) clean() (count int) {
	now := time.Now()

	for key, value := range v.data {
		if now.Sub(value.Timestamp) > v.timeToLive {
			delete(v.data, key)
			count += 1
		}
	}
	return
}

func (v *Volatile[K, V]) cleanupRoutine(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		v.clean()
	}
}

func NewVolatile[K comparable, V any](timeToLive time.Duration, cleanupInterval time.Duration) *Volatile[K, V] {
	v := &Volatile[K, V]{
		timeToLive: timeToLive,
	}
	v.Clear()
	go v.cleanupRoutine(cleanupInterval)
	return v
}

func (v *Volatile[K, V]) Has(key K) (ok bool) {
	v.clean()
	_, ok = v.data[key]
	return
}

func (v *Volatile[K, V]) Get(key K) (*V, error) {
	v.clean()
	element, ok := v.data[key]

	if !ok {
		return nil, fmt.Errorf("not found")
	}

	return element.value, nil
}

func (v *Volatile[K, V]) Remove(key K) (*V, error) {
	v.clean()
	value, ok := v.data[key]

	if !ok {
		return nil, fmt.Errorf("not found")
	}

	delete(v.data, key)
	return value.value, nil
}

func (v *Volatile[K, V]) Length() int {
	v.clean()
	return len(v.data)
}

func (v *Volatile[K, V]) Clear() {
	v.data = make(map[K]Element[V])
}

func (v *Volatile[K, V]) Set(key K, value *V) {
	v.data[key] = Element[V]{value: value, Timestamp: time.Now()}
	v.clean()
}
