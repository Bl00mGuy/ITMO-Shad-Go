//go:build !solution

package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	wordCounter := make(map[string]int)

	for _, file := range os.Args[1:] {
		file, err := os.Open(file)
		check(err)
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			word := strings.TrimSpace(scanner.Text())
			wordCounter[word]++
		}
	}

	for word, count := range wordCounter {
		if count > 1 {
			fmt.Printf("%d\t%s\n", count, word)
		}
	}
}
