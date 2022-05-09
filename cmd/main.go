package main

import (
	"fmt"
	"sync"

	"github.com/go-ping/ping"
	"github.com/jlineaweaver/wof"
)

var wg sync.WaitGroup

func main() {

	websites := wof.Scan("websites")

	success := make(chan string, 100)
	fail := make(chan string, 100)

	go TestWebsites(websites, success, fail)

  agg := wof.Aggregator{}
	Aggregate(&agg, success, fail)
	agg.Results()
}

func Aggregate(agg *wof.Aggregator, suc <-chan string, fail <-chan string) {
	for {
		select {
		case url, ok := <-suc:
			if !ok {
				fmt.Println(url)
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
}

func TestWebsites(websites []string, success chan<- string, fail chan<- string) {
	defer close(success)
	defer close(fail)

	wg := sync.WaitGroup{}
	for _, website := range websites {
		wg.Add(1)
		go TestWebsite(&wg, website, success, fail)
	}

	wg.Wait()
}

func TestWebsite(wg *sync.WaitGroup, website string, success chan<- string, fail chan<- string) {
	pinger, err := ping.NewPinger(website)
	pinger.SetPrivileged(true)

	if err != nil {
		fail <- website
		wg.Done()
		return
	}
	pinger.Count = 3
	pinger.Run() // blocks until finished
	success <- website
	wg.Done()
}
