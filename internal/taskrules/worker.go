package taskrules

import "sync"

func RunConcurrentAdds(counter *Counter, goroutines, iterations int) {
	var wg sync.WaitGroup
	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				counter.add(1)
			}
		}()
	}
	wg.Wait()
}
