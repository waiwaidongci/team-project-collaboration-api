package taskrules

import "testing"

func TestCounterConcurrentAdds(t *testing.T) {
	counter := NewCounter()
	RunConcurrentAdds(counter, 8, 1000)
	if got, want := counter.Value(), 8*1000; got != want {
		t.Fatalf("counter value = %d, want %d", got, want)
	}
}
