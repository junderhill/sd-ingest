package util

import "sync"

type SafeInt64Count struct {
	mu    sync.Mutex
	count int64
}

type SafeIntCount struct {
	mu    sync.Mutex
	count int
}

func (c *SafeIntCount) Increment(value int) {
	c.mu.Lock()
	c.count += value
	c.mu.Unlock()
}

func (c *SafeInt64Count) Increment(value int64) {
	c.mu.Lock()
	c.count += value
	c.mu.Unlock()
}

func (c *SafeIntCount) Value() int {
	return c.count
}

func (c *SafeInt64Count) Value() int64 {
	return c.count
}
