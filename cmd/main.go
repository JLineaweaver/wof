package main

import (
	"fmt"
	"strings"
	"sync"

	"github.com/go-ping/ping"
)

var wg sync.WaitGroup

type Aggregator struct {
	lock      sync.Mutex
	successes []string
	failures  []string
}

func (agg *Aggregator) Success(url string) {
	defer agg.lock.Unlock()
	agg.lock.Lock()
	agg.successes = append(agg.successes, url)
}

func (agg *Aggregator) Failure(url string) {
	defer agg.lock.Unlock()
	agg.lock.Lock()
	agg.failures = append(agg.failures, url)
}

func (agg *Aggregator) Results() {
	fmt.Printf("Number of successes: %d\n", len(agg.successes))
	fmt.Printf("Number of failures: %d\n", len(agg.failures))
	if len(agg.failures) > 0 {
		fmt.Printf("Failures: %s\n", strings.Join(agg.failures, ", "))
	}
}

func main() {

  worker := make(chan string, 10)

 // reader := wof.Reader{worker}

	suc := make(chan string, 10)
	fail := make(chan string, 10)

	agg := Aggregate(suc, fail)

	//wg.Wait()
	agg.Results()
}

func Aggregate(suc <-chan string, fail <-chan string) Aggregator {
	agg := Aggregator{}
	for i := 0; i < 4; i++ {
		select {
		case url := <-suc:
			agg.Success(url)
		case url := <-fail:
			agg.Failure(url)
		}
	}

	return agg
}

func TestWebsite(website string, suc chan<- string, fail chan<- string) {
	pinger, err := ping.NewPinger(website)
	pinger.SetPrivileged(true)

	if err != nil {
		fail <- website
	}
	pinger.Count = 3
	pinger.Run() // blocks until finished

	suc <- website
}
