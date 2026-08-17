package taskrules

func NewCounter() *Counter {
	return &Counter{}
}

func (c *Counter) Reset() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.value = 0
}

func (c *Counter) Clone() *Counter {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return &Counter{value: c.value}
}
