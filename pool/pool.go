package pool

import (
	"errors"
	"log"
	"runtime/debug"
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

	tasks chan func()
	wg    sync.WaitGroup

	stopOnce sync.Once
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
		tasks:        make(chan func(), queueSize),
	}

	// Start workers (Step 2: minimal happy path)
	for i := 0; i < p.workersCount; i++ {
		p.wg.Add(1)
		go func() {
			defer p.wg.Done()
			for task := range p.tasks {
				if task == nil {
					continue
				}
				func() {
					defer func() {
						if r := recover(); r != nil {
							log.Printf("pool: task panic: %v\n%s", r, string(debug.Stack()))
						}
						if hook := p.cfg.onTaskDone; hook != nil {
							defer func() { _ = recover() }()
							hook()
						}
					}()
					task()
				}()
			}
		}()
	}
	return p, nil
}

func (p *poolImpl) Submit(task func()) error {
	p.mu.Lock()
	if p.stopped {
		p.mu.Unlock()
		return ErrStopped
	}
	// protected send under lock to avoid racing with close(p.tasks)
	select {
	case p.tasks <- task:
		p.mu.Unlock()
		return nil
	default:
		p.mu.Unlock()
		return ErrQueueFull
	}
}

func (p *poolImpl) Stop() error {
	p.stopOnce.Do(func() {
		p.mu.Lock()
		if !p.stopped {
			p.stopped = true
			close(p.tasks)
		}
		p.mu.Unlock()
		// wait for workers to exit
		p.wg.Wait()
	})
	return nil
}
