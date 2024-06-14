package volatile

import (
	"testing"
	"time"
)

func TestVolatile_SetGet(t *testing.T) {
	cache := NewVolatile[string, string](2*time.Second, 1*time.Second)

	key := "key1"
	value := "value1"
	cache.Set(key, &value)

	got, err := cache.Get(key)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if *got != value {
		t.Errorf("got %v, want %v", *got, value)
	}

	if !cache.Has(key) {
		t.Errorf("expected key %v to exist", key)
	}
}

func TestVolatile_Remove(t *testing.T) {
	cache := NewVolatile[string, string](2*time.Second, 1*time.Second)

	key := "key1"
	value := "value1"
	cache.Set(key, &value)

	got, err := cache.Remove(key)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if *got != value {
		t.Errorf("got %v, want %v", *got, value)
	}

	if cache.Has(key) {
		t.Errorf("expected key %v to be removed", key)
	}
}

func TestVolatile_Length(t *testing.T) {
	cache := NewVolatile[string, int](100*time.Millisecond, 50*time.Millisecond)

	if cache.Length() != 0 {
		t.Errorf("expected length to be 0, got %d", cache.Length())
	}

	v1 := 1
	v2 := 2

	cache.Set("key1", &v1)
	cache.Set("key2", &v2)

	if cache.Length() != 2 {
		t.Errorf("expected length to be 2, got %d", cache.Length())
	}

	// Wait for the elements to expire and be cleaned up
	time.Sleep(300 * time.Millisecond)

	if cache.Length() != 0 {
		t.Errorf("expected length to be 0 after expiration, got %d", cache.Length())
	}
}

func TestVolatile_Clear(t *testing.T) {
	cache := NewVolatile[string, int](2*time.Second, 1*time.Second)

	v1 := 1
	v2 := 2
	v3 := 3

	// Insert multiple items
	cache.Set("key1", &v1)
	cache.Set("key2", &v2)
	cache.Set("key3", &v3)

	// Ensure items are set
	if !cache.Has("key1") || !cache.Has("key2") || !cache.Has("key3") {
		t.Fatal("expected keys to be set")
	}

	// Clear the cache
	cache.Clear()

	// Ensure all items are cleared
	if cache.Has("key1") || cache.Has("key2") || cache.Has("key3") {
		t.Fatal("expected all keys to be cleared")
	}

	// Ensure map is empty
	if len(cache.data) != 0 {
		t.Errorf("expected map to be empty, got %d elements", len(cache.data))
	}
}

func TestVolatile_Clean(t *testing.T) {
	cache := NewVolatile[string, string](100*time.Millisecond, 50*time.Millisecond)

	key := "key1"
	value := "value1"
	cache.Set(key, &value)

	time.Sleep(150 * time.Millisecond) // Wait for the element to expire

	if cache.Has(key) {
		t.Errorf("expected key %v to be expired and removed", key)
	}

	_, err := cache.Get(key)
	if err == nil {
		t.Errorf("expected error when getting expired key %v", key)
	}
}

func TestVolatile_AutomaticCleanup(t *testing.T) {
	cache := NewVolatile[string, string](100*time.Millisecond, 50*time.Millisecond)

	key1 := "key1"
	value1 := "value1"
	cache.Set(key1, &value1)

	key2 := "key2"
	value2 := "value2"
	cache.Set(key2, &value2)

	time.Sleep(150 * time.Millisecond) // Wait for the elements to expire

	if cache.Has(key1) || cache.Has(key2) {
		t.Errorf("expected all keys to be expired and removed")
	}
}
