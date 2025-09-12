package pool_test

import (
    "errors"
    pool "example.com/ksb/pool/pool"
    "sync"
    "sync/atomic"
    "testing"
)

func TestNewPool_InvalidConfig(t *testing.T) {
    cases := [][2]int{
        {-1, 0},
        {0, -1},
        {-1, -1},
    }
    for _, c := range cases {
        if p, err := pool.NewPool(c[0], c[1]); err == nil || !errors.Is(err, pool.ErrInvalidConfig) || p != nil {
            t.Fatalf("expected ErrInvalidConfig for workers=%d queue=%d, got p=%v err=%v", c[0], c[1], p, err)
        }
    }
}

func TestNewPool_ValidConfig(t *testing.T) {
    hooks := []func(){nil, func() {}}
    cases := [][2]int{
        {0, 0},
        {1, 0},
        {0, 1},
        {2, 3},
    }
    for _, c := range cases {
        for _, h := range hooks {
            p, err := pool.NewPool(c[0], c[1], pool.WithOnTaskDone(h))
            if err != nil || p == nil {
                t.Fatalf("unexpected error for workers=%d queue=%d: %v", c[0], c[1], err)
            }
        }
    }
}

func TestStop_NoWorkers_NoRace(t *testing.T) {
    p, err := pool.NewPool(0, 0)
    if err != nil {
        t.Fatalf("new pool: %v", err)
    }

    const goroutines = 32
    const iterations = 64
    var wg sync.WaitGroup
    var errCount int64
    wg.Add(goroutines)
    for g := 0; g < goroutines; g++ {
        go func() {
            defer wg.Done()
            for i := 0; i < iterations; i++ {
                if err := p.Stop(); err != nil {
                    atomic.AddInt64(&errCount, 1)
                }
            }
        }()
    }
    wg.Wait()
    if atomic.LoadInt64(&errCount) != 0 {
        t.Fatalf("stop returned errors, count=%d", errCount)
    }
}


