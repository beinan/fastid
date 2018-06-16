package fastid

import (
	"sync"
	"testing"
	"time"
)

func TestGenID(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(100)                        //using 100 goroutine to generate 10000 ids
	results := make(chan int64, 10000) //store result
	for i := 0; i < 100; i++ {
		go func() {
			for i := 0; i < 100; i++ {
				id := GenInt64ID()
				t.Logf("id: %b \t %x \t %d", id, id, id)
				results <- id
			}
			wg.Done()
		}()
	}
	wg.Wait()

	m := make(map[int64]bool)
	for i := 0; i < 10000; i++ {
		select {
		case id := <-results:
			if _, ok := m[id]; ok {
				t.Errorf("Found duplicated id: %x", id)
				return
			} else {
				m[id] = true
			}
		case <-time.After(2 * time.Second):
			t.Errorf("Expect 10000 ids in results, but got %d", i)
			return
		}
	}
}

func BenchmarkGenID(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GenInt64ID()
	}
}

func BenchmarkGenIDP(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			GenInt64ID()
		}
	})
}
