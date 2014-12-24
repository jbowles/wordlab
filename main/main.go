package main

import (
	"fmt"
	"github.com/jbowles/wordlab"
)

func main() {
	fp := "wordlab_test.csv"
	category := "Testing"
	wordlab.CsvCreateFileWithHeaders(false, fp, []string{"ByteRange", "FitValue", "Word", "Category"})

	fmt.Printf("Encode: %v\n", wordlab.Format("this", category, fp))
	fmt.Printf("Encode: %v\n", wordlab.Format("ttis", category, fp))
	fmt.Printf("Encode: %v\n", wordlab.Format("Ω≈ç∂´", category, fp))
	fmt.Printf("Encode: %v\n", wordlab.Format("Ω≈ç∂´ ", category, fp))
}
