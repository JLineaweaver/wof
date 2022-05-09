package main

import (
	"sync"

	"github.com/go-ping/ping"
	"github.com/jlineaweaver/wof"
)

var wg sync.WaitGroup

func main() {

	worker := make(chan string, 100)
	go wof.Scan("websites", worker)

	success := make(chan string, 100)
	fail := make(chan string, 100)

	go TestWebsites(worker, success, fail)

	agg := Aggregate(success, fail)
	agg.Results()
}

func Aggregate(suc <-chan string, fail <-chan string) wof.Aggregator {
	agg := wof.Aggregator{}
	for {
		select {
		case url, ok := <-suc:
			if !ok {
				suc = nil
				break
			}
			agg.Success(url)
		case url, ok := <-fail:
			if !ok {
				fail = nil
				break
			}
			agg.Failure(url)
		}
		if suc == nil && fail == nil {
			break
		}
	}
	return agg
}

func TestWebsites(websites <-chan string, success chan<- string, fail chan<- string) {
	defer close(success)
	defer close(fail)

	for {
		select {
		case website, ok := <-websites:
			if !ok {
				return
			}
			pinger, err := ping.NewPinger(website)
			pinger.SetPrivileged(true)

			if err != nil {
				fail <- website
				break
			}
			pinger.Count = 3
			pinger.Run() // blocks until finished
			success <- website
		}
	}
}
