package pool_test

import (
	"sync/atomic"
	"testing"

	pool "github.com/OptonGroup/kaspersky-safeboard-go-sandbox-development-team/pool"
)

func TestTaskPanic_DoesNotKillWorker_AndOnTaskDoneCalled(t *testing.T) {
	var doneCount int32
	p, err := pool.NewPool(2, 4, pool.WithOnTaskDone(func() { atomic.AddInt32(&doneCount, 1) }))
	if err != nil {
		t.Fatalf("new pool: %v", err)
	}

	// Упавшая задача + нормальная задача после неё должны обе отработать
	if err := p.Submit(func() { panic("boom") }); err != nil {
		t.Fatalf("submit panic task: %v", err)
	}
	ran := int32(0)
	if err := p.Submit(func() { atomic.AddInt32(&ran, 1) }); err != nil {
		t.Fatalf("submit normal task: %v", err)
	}

	if err := p.Stop(); err != nil {
		t.Fatalf("stop: %v", err)
	}

	if atomic.LoadInt32(&ran) != 1 {
		t.Fatalf("expected normal task to run")
	}
	if atomic.LoadInt32(&doneCount) != 2 {
		t.Fatalf("expected onTaskDone called twice, got %d", doneCount)
	}
}

