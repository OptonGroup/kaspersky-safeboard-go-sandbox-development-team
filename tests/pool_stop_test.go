package pool_test

import (
	"sync/atomic"
	"testing"
	"time"

	pool "example.com/ksb/pool/pool"
)

func TestStop_RejectsNewSubmits(t *testing.T) {
	p, err := pool.NewPool(2, 4)
	if err != nil {
		t.Fatalf("new pool: %v", err)
	}
	if err := p.Stop(); err != nil {
		t.Fatalf("stop: %v", err)
	}
	if err := p.Submit(func() {}); err != pool.ErrStopped {
		t.Fatalf("expected ErrStopped after Stop, got %v", err)
	}
}

func TestStop_WaitsForTasks(t *testing.T) {
	p, err := pool.NewPool(2, 4)
	if err != nil {
		t.Fatalf("new pool: %v", err)
	}

	var ran int32
	// Две долгие задачи должны завершиться до выхода Stop
	if err := p.Submit(func() { time.Sleep(80 * time.Millisecond); atomic.AddInt32(&ran, 1) }); err != nil {
		t.Fatalf("submit: %v", err)
	}
	if err := p.Submit(func() { time.Sleep(80 * time.Millisecond); atomic.AddInt32(&ran, 1) }); err != nil {
		t.Fatalf("submit: %v", err)
	}

	start := time.Now()
	if err := p.Stop(); err != nil {
		t.Fatalf("stop: %v", err)
	}
	elapsed := time.Since(start)
	if atomic.LoadInt32(&ran) != 2 {
		t.Fatalf("expected all tasks to finish, ran=%d", ran)
	}
	if elapsed < 80*time.Millisecond {
		t.Fatalf("stop returned too early: %v", elapsed)
	}
}

func TestStop_IsIdempotent(t *testing.T) {
	p, err := pool.NewPool(1, 1)
	if err != nil {
		t.Fatalf("new pool: %v", err)
	}
	if err := p.Stop(); err != nil {
		t.Fatalf("stop: %v", err)
	}
	// Повторные вызовы не должны паниковать/возвращать ошибку
	if err := p.Stop(); err != nil {
		t.Fatalf("stop again: %v", err)
	}
}
