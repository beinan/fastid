package fastid

import (
	"testing"
	//	"time"
)

func TestGenID(t *testing.T) {
	for i := 0; i < 0; i++ {
		//t.Log("generating id")
		go func() {
			for i := 1; i < 20; i++ {
				id := GenInt64ID()
				t.Logf("id: %b \t %x \t %d", id, id, id)
			}
		}()
	}
	//time.Sleep(5 * time.Second)
}

func TestIP(t *testing.T) {
	ip, err := getIP()
	t.Logf("ip: %v %v", ip, err)
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
