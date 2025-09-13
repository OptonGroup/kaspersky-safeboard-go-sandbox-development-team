package pool_test

import (
	pool "github.com/OptonGroup/kaspersky-safeboard-go-sandbox-development-team/pool"
)

// This example demonstrates basic usage of the worker pool.
func Example() {
	p, _ := pool.NewPool(2, 4)
	defer p.Stop()

	p.Submit(func() { /* do work */ })
	p.Submit(func() { /* do more work */ })
	// Output:
}
