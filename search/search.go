package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	// 1. Read input from user
	var input string
	fmt.Print("Enter search word: ")
	fmt.Scanln(&input)

	path := filepath.Join("output", input+".txt")
	// 2. Open file with that name
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()
	fmt.Println()
	// 3. Read content and list URLs
	urlSet := make(map[string]bool) // to remove duplicates
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		url := scanner.Text()
		if !urlSet[url] {
			fmt.Println("https://" + url)
			urlSet[url] = true
		}
	}
	fmt.Println()
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
	}
}
