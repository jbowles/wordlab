package wordlab

import (
	"bufio"
	"fmt"

	"os"
)

type StopWordList struct {
	IsStopWord map[string]bool
	Total      int
}

func StopWords(file_name string) *StopWordList {
	var total_count int
	swb := &StopWordList{
		IsStopWord: make(map[string]bool),
	}
	file, ferr := os.Open(file_name)

	if ferr != nil {
		fmt.Println("Error reading file: ", ferr)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		swb.IsStopWord[scanner.Text()] = true
		total_count += 1
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("**error:", err)
	}
	swb.Total = total_count
	return swb
}
