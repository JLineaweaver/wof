package wof

import (
	"fmt"
	"strings"
	"sync"
)

type Aggregator struct {
	lock      sync.Mutex
	successes []string
	failures  []string
}

func (agg *Aggregator) Success(url string) {
	defer agg.lock.Unlock()
	agg.lock.Lock()
	agg.successes = append(agg.successes, url)
	fmt.Printf("Success adding %s\n", url)
}

func (agg *Aggregator) Failure(url string) {
	defer agg.lock.Unlock()
	agg.lock.Lock()
	agg.failures = append(agg.failures, url)
	fmt.Printf("Failure adding %s\n", url)
}

func (agg *Aggregator) Results() {
	fmt.Printf("Number of successes: %d\n", len(agg.successes))
	fmt.Printf("Number of failures: %d\n", len(agg.failures))
	if len(agg.failures) > 0 {
		fmt.Printf("Failures: %s\n", strings.Join(agg.failures, ", "))
	}
}


