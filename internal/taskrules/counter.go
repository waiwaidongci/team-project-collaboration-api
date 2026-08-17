package taskrules

import "sync"

type Counter struct {
	mu    sync.RWMutex
	value int
}

func (c *Counter) add(n int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.value += n
}
