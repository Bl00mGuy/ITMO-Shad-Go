package main

import (
	"bufio"
	"os"
)

func readInput(input *os.File) chan string {
	ch := make(chan string)
	go func() {
		scanner := bufio.NewScanner(input)
		for scanner.Scan() {
			ch <- scanner.Text()
		}
		close(ch)
	}()
	return ch
}
