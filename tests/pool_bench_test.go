package pool_test

import (
	"testing"

	pool "github.com/OptonGroup/kaspersky-safeboard-go-sandbox-development-team/pool"
)

func benchmarkPool(b *testing.B, workers, queue int) {
	p, err := pool.NewPool(workers, queue)
	if err != nil {
		b.Fatalf("new pool: %v", err)
	}
	defer p.Stop()

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := p.Submit(func() {}); err != nil {
			// If queue is full or stopped during benchmark end, just continue
			// to avoid skewing results with sleeps.
		}
	}
}

func BenchmarkPool_W1_Q1(b *testing.B)  { benchmarkPool(b, 1, 1) }
func BenchmarkPool_W4_Q16(b *testing.B) { benchmarkPool(b, 4, 16) }
func BenchmarkPool_W8_Q64(b *testing.B) { benchmarkPool(b, 8, 64) }
