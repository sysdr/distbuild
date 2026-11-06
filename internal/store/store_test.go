package store

import (
	"fmt"
	"sync"
	"testing"
)

func TestSetGetDelete(t *testing.T) {
	s := New()
	if _, ok := s.Get("a"); ok {
		t.Fatal("expected missing key")
	}
	s.Set("a", "1")
	if v, ok := s.Get("a"); !ok || v != "1" {
		t.Fatalf("got %q,%v want 1,true", v, ok)
	}
	if !s.Delete("a") {
		t.Fatal("delete should report existing key")
	}
	if s.Delete("a") {
		t.Fatal("second delete should report missing key")
	}
}

// Run with: go test -race ./...
func TestConcurrentAccess(t *testing.T) {
	s := New()
	var wg sync.WaitGroup
	for w := 0; w < 50; w++ {
		wg.Add(1)
		go func(w int) {
			defer wg.Done()
			for i := 0; i < 1000; i++ {
				k := fmt.Sprintf("k%d-%d", w, i)
				s.Set(k, "v")
				s.Get(k)
			}
		}(w)
	}
	wg.Wait()
	if got := s.Len(); got != 50*1000 {
		t.Fatalf("len=%d want %d", got, 50*1000)
	}
}
