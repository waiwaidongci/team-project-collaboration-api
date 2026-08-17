package taskrules

func NewCounter() *Counter {
	return &Counter{}
}

func (c *Counter) Reset() {
	c.value = 0
}

func (c *Counter) Clone() *Counter {
	return &Counter{value: c.value}
}
