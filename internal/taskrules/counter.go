package taskrules

type Counter struct {
	value int
}

func (c *Counter) add(n int) {
	c.value += n
}
