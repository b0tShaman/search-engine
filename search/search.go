package main

import (
	"encoding/gob"
	"fmt"
	"os"
	"sort"
)

const (
	MAX_RESULTS = 10
)

func main() {
	// Load the inverted index from the file
	var invertedIndex map[string]map[string]int

	file, _ := os.Open("index.gob")
	defer file.Close()
	gob.NewDecoder(file).Decode(&invertedIndex)

	// Prompt user for input
	var input string
	fmt.Print("Enter search word: ")
	fmt.Scanln(&input)

	fmt.Println()
	// Retrieve and sort URLs based on frequency
	urlSet, exists := invertedIndex[input]
	if !exists {
		fmt.Println("No results found for the given word.")
		return
	}
	// Convert map to slice of key-value pairs
	type kv struct {
		Key   string
		Value int
	}

	var pairs []kv
	for k, v := range urlSet {
		pairs = append(pairs, kv{k, v})
	}

	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].Value > pairs[j].Value
	})
	// Print sorted URLs with counts only top 10
	for i, p := range pairs {
		if i >= MAX_RESULTS {
			break
		}
		fmt.Printf("https://%s -> %d\n", p.Key, p.Value)
	}
}
