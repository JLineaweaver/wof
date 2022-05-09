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

	agg := Aggregate(success, fail)
	agg.Results()
}

func Aggregate(suc <-chan string, fail <-chan string) wof.Aggregator {
	agg := wof.Aggregator{}
	for {
		select {
		case url, ok := <-suc:
			if !ok {
				fmt.Println("closed suc")
				suc = nil
				break
			}
			fmt.Println("hit suc")
			agg.Success(url)
		case url, ok := <-fail:
			if !ok {
				fmt.Println("closed fail")
				fail = nil
				break
			}
			fmt.Println("hit fail")
			agg.Failure(url)
		}
		if suc == nil && fail == nil {
			break
		}
	}
	return agg
}

func TestWebsites(websites []string, success chan<- string, fail chan<- string) {
	defer close(success)
	defer close(fail)

  for _, website := range websites{
			go TestWebsite(website, success, fail)
	}
}

func TestWebsite(website string, success chan<- string, fail chan<- string) {
	fmt.Printf("We testing %s\n", website)
	pinger, err := ping.NewPinger(website)
	pinger.SetPrivileged(true)

	if err != nil {
		fail <- website
		return
	}
	pinger.Count = 3
	pinger.Run() // blocks until finished
	success <- website
}
