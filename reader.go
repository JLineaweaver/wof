package wof

import (
	"bufio"
	"os"
)

func Scan(file string) []string {
	results := []string{}

	f, err := os.Open("websites")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		results = append(results, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}
	return results
}
