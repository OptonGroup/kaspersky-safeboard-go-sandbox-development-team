package pool

import (
	"errors"
	"sync"
)

// Public errors
var (
	ErrQueueFull     = errors.New("pool: queue is full")
	ErrStopped       = errors.New("pool: stopped")
	ErrInvalidConfig = errors.New("pool: invalid config")
)

// Pool describes a minimal worker pool API
type Pool interface {
	// Submit schedules a task for execution.
	// In later steps it will enqueue or return ErrQueueFull when the queue is full.
	Submit(task func()) error

	// Stop gracefully stops the pool. In later steps it will become idempotent
	// and wait for running tasks to finish.
	Stop() error
}

// Option configures Pool behavior
type Option func(*config)

type config struct {
	onTaskDone func()
}

// WithOnTaskDone sets a hook that is called when a task finishes.
// The hook may be nil. It will be invoked in later steps.
func WithOnTaskDone(hook func()) Option {
	return func(c *config) {
		c.onTaskDone = hook
	}
}

// poolImpl is a minimal skeleton implementation for Step 1
type poolImpl struct {
	workersCount int
	queueSize    int

	cfg config

	mu      sync.Mutex
	stopped bool
}

// NewPool constructs a new Pool instance.
// Validates workers and queueSize are non-negative. Zero values are allowed.
func NewPool(workers, queueSize int, opts ...Option) (Pool, error) {
	if workers < 0 || queueSize < 0 {
		return nil, ErrInvalidConfig
	}

	cfg := config{}
	for _, opt := range opts {
		if opt != nil {
			opt(&cfg)
		}
	}

	p := &poolImpl{
		workersCount: workers,
		queueSize:    queueSize,
		cfg:          cfg,
	}
	return p, nil
}

// Submit: Step 1 returns ErrStopped until workers/queue are implemented in Step 2.
func (p *poolImpl) Submit(task func()) error {
	p.mu.Lock()
	stopped := p.stopped
	p.mu.Unlock()
	if stopped {
		return ErrStopped
	}
	return ErrStopped
}

// Stop marks the pool as stopped. Idempotency and waiting logic will be added later.
func (p *poolImpl) Stop() error {
	p.mu.Lock()
	if p.stopped {
		p.mu.Unlock()
		return nil
	}
	p.stopped = true
	p.mu.Unlock()
	return nil
}
