package wof

import (
	"bufio"
	"os"
)

func Scan(file string, c chan<- string) {
	defer close(c)

	f, err := os.Open("websites")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		c <- scanner.Text()
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}
}
