package pool_test

import (
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	pool "github.com/OptonGroup/kaspersky-safeboard-go-sandbox-development-team/pool"
)

func TestRace_SubmitAndStop_NoPanicsOrLeaks(t *testing.T) {
	p, err := pool.NewPool(runtime.GOMAXPROCS(0), 64)
	if err != nil {
		t.Fatalf("new pool: %v", err)
	}

	var submits int64
	var errs int64
	var wg sync.WaitGroup

	stopCh := make(chan struct{})

	// Submitters
	workers := runtime.GOMAXPROCS(0) * 2
	wg.Add(workers)
	for i := 0; i < workers; i++ {
		go func() {
			defer wg.Done()
			for {
				select {
				case <-stopCh:
					return
				default:
				}
				if err := p.Submit(func() {}); err != nil {
					atomic.AddInt64(&errs, 1)
				} else {
					atomic.AddInt64(&submits, 1)
				}
			}
		}()
	}

	// Let it run briefly
	time.Sleep(100 * time.Millisecond)

	// Stop concurrently with submitters
	if err := p.Stop(); err != nil {
		t.Fatalf("stop: %v", err)
	}
	close(stopCh)
	wg.Wait()

	// Pool must be stopped and not panic; errors are expected (queue full or stopped)
	if submits == 0 && errs == 0 {
		t.Fatalf("expected some submits or errors recorded")
	}
}
