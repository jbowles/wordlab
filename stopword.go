package wordlab

import (
	"bufio"
	"fmt"

	"os"
)

type StopWordBucket struct {
	IsStopWord map[string]bool
	Total      int
}

func StopWords(file_name string) *StopWordBucket {
	var total_count int
	swb := &StopWordBucket{
		IsStopWord: make(map[string]bool),
	}
	//file, err := os.Open("README.md")
	file, ferr := os.Open(file_name)

	if ferr != nil {
		fmt.Println("Error reading file: ", ferr)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		//swb.Words = append(swb.Words, []string{scanner.Text(), ""})
		swb.IsStopWord[scanner.Text()] = true
		total_count += 1
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("**error:", err)
	}
	swb.Total = total_count
	return swb
}
