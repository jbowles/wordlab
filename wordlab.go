/*
* wordlab package pre-processes and formats words for further processing, specifically for classification or clustering algorithms (knn, k-means, x-means, etc...).
* It creates a unique floating point numeral for each unique word and writes to a file.
 */
package wordlab

import (
	"encoding/csv"
	"fmt"
	"os"
)

const (
	Eqlzr = 0.13
)

var (
	posTotal float64
	chrTotal float64
)

type WordBucket struct {
	Bucket   []BytePosChar
	Word     string
	FitValue float64
	Category string
}

type BytePosChar struct {
	BytePosition  float64
	ByteCharacter float64
}

func Format(word, category, file_path string) *WordBucket {
	bucket := &WordBucket{
		Word:     word,
		Category: category,
	}

	for pos, chr := range []byte(word) {
		bucket.Bucket = append(bucket.Bucket, setBytePosChar(float64(pos), float64(chr)))
	}
	bucket.setFitValue()
	bucket.csvBucketWriter(file_path)
	return bucket
}

// CsvCreatFileWithHeaders creates csv with headers.
// This is a destructive function as it will overwrite existing files with the same name!
func CsvCreateFileWithHeaders(force bool, file_path string, headers []string) {
	file_exists := fileExist(file_path)
	switch file_exists {
	case false:
		if force == false {
			return
		}
	}
	csvfile, err := os.Create(file_path)
	if err != nil {
		fmt.Println(err)
		return
	}
	writer := csv.NewWriter(csvfile)
	writer.Write(headers)
	writer.Flush()
	csvfile.Close()
}

func setBytePosChar(position, character float64) BytePosChar {
	return BytePosChar{
		BytePosition:  position,
		ByteCharacter: character,
	}
}

func (wb *WordBucket) setFitValue() {
	for _, bpc := range wb.Bucket {
		chrTotal += bpc.ByteCharacter
		posTotal += bpc.BytePosition
	}
	wb.FitValue = ((posTotal * Eqlzr) / chrTotal)
}

func fileExist(file_path string) bool {
	_, err := os.OpenFile(file_path, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0660)
	return err == nil
}

// CsvBucketWriter will append or create file to write out a WordBucket.
func (wb *WordBucket) csvBucketWriter(file_path string) {
	csvfile, err := os.OpenFile(file_path, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0660)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer csvfile.Close()

	writer := csv.NewWriter(csvfile)

	fit_value := fmt.Sprintf("%G", wb.FitValue)
	bucket := fmt.Sprintf("%v", wb.Bucket)
	writeErr := writer.Write(
		[]string{
			bucket,
			fit_value,
			wb.Word,
			wb.Category,
		},
	)

	if writeErr != nil {
		fmt.Println(writeErr)
	}
	writer.Flush()
}
