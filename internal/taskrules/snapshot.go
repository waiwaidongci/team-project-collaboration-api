package taskrules

func (c *Counter) Value() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.value
}
