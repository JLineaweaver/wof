package wof

import (
	"bufio"
	"log"
	"os"
)

type Reader struct {
	c chan<- string
}

func (r *Reader) Scan(file string) error {
	f, err := os.Open("websites")
	if err != nil {
		return err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	// optionally, resize scanner's capacity for lines over 64K, see next example

	for scanner.Scan() {
		//wg.Add(1)
		r.c <- scanner.Text()
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return nil
}
