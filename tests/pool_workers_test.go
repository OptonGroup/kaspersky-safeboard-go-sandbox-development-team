package pool_test

import (
	"testing"
	"time"

	pool "github.com/OptonGroup/kaspersky-safeboard-go-sandbox-development-team/pool"
)

func TestWorkers_RunTasksInParallel(t *testing.T) {
	workers := 4
	queue := 8
	p, err := pool.NewPool(workers, queue)
	if err != nil {
		t.Fatalf("new pool: %v", err)
	}
	defer p.Stop()

	taskDelay := 50 * time.Millisecond
	tasks := 8

	start := time.Now()
	for i := 0; i < tasks; i++ {
		if err := p.Submit(func() { time.Sleep(taskDelay) }); err != nil {
			t.Fatalf("submit: %v", err)
		}
	}

	// В этом шаге Stop не ждёт, поэтому подождём верхнюю оценку времени, соответствующую параллельности.
	// Последовательное время ~ tasks*taskDelay = 400ms.
	// Ожидаем, что параллельно на 4 воркерах укладываемся существенно меньше, проверим < 300ms.
	time.Sleep(300 * time.Millisecond)

	elapsed := time.Since(start)
	if elapsed >= 400*time.Millisecond {
		t.Fatalf("expected parallel speedup, elapsed=%v", elapsed)
	}
}

func TestSubmit_QueueLimit_ReturnsErrQueueFull(t *testing.T) {
	p, err := pool.NewPool(0, 2)
	if err != nil {
		t.Fatalf("new pool: %v", err)
	}
	defer p.Stop()

	if err := p.Submit(func() {}); err != nil {
		t.Fatalf("submit 1: %v", err)
	}
	if err := p.Submit(func() {}); err != nil {
		t.Fatalf("submit 2: %v", err)
	}
	if err := p.Submit(func() {}); err == nil {
		t.Fatalf("expected ErrQueueFull, got nil")
	} else if err != pool.ErrQueueFull {
		t.Fatalf("expected ErrQueueFull, got %v", err)
	}
}
