package wordlab

import (
	"encoding/csv"
	"fmt"
	tkz "github.com/jbowles/nlpt_tkz"
	"io"
	"os"
	"strconv"
	"strings"
)

func fileExist(file_path string) bool {
	_, err := os.OpenFile(file_path, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	return err == nil
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

func FormatCSV(word, category, write_file, format string) *WordBucket {
	bucket := &WordBucket{
		Word:     word,
		Category: category,
	}

	for pos, chr := range []byte(word) {
		bucket.Bucket = append(bucket.Bucket, setBytePosChar(float64(pos), float64(chr)))
	}
	bucket.setFitValue()
	switch format {
	case "amit":
		bucket.csvBucketKnnAmitWriter(write_file)
	case "normal":
		bucket.csvBucketKnnWriter(write_file)
	default:
		bucket.csvBucketKnnWriter(write_file)
	}

	return bucket
}

func ReadCsvFormatCsv(read_file, write_file, category, format string) {
	csvFile, err := os.Open(read_file)
	stop := StopWords("datasets/stopwords/stopwords.txt")
	defer csvFile.Close()
	if err != nil {
		panic(err)
	}
	csvReader := csv.NewReader(csvFile)
	csvReader.TrimLeadingSpace = true
	csvReader.FieldsPerRecord = 1
	for {
		fields, err := csvReader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
		for _, row := range fields {
			// WRITE WHOLE ROW
			//FormatCSV(row, category, write_file)
			tokens, _ := tkz.Tokenize("bukt", strings.ToLower(row))
			for _, token := range tokens {
				if !stop.IsStopWord[token] {
					FormatCSV(token, category, write_file, format)
				}
			}
		}
	}
}

// ConcatRuneSlice concatenates two slices of runes.
// Takes two arguments that must be slices of type Rune, order of args is not important.
// Returns a new slice which is the result of copying the 2 slices passed in a args.
//
//  slice_one := []rune("this, that, other")
//  slice_two := []rune("when, where, why")
//  goko.ConcatRuneSlice(slice_one,slice_two)
//  //new_slice =>
func ConcatStringSlice(slice1, slice2 []string) []string {
	new_slice := make([]string, len(slice1)+len(slice2))
	copy(new_slice, slice1)
	copy(new_slice[len(slice1):], slice2)
	return new_slice
}

func (wb *WordBucket) BytePosCharToString() []string {
	size := len(wb.Bucket)
	slice0 := []string{}
	tmp := []string{}
	for _, i := range wb.Bucket {
		tmp = append(tmp, strconv.Itoa(int(i.ByteCharacter)))
	}
	slice1 := ConcatStringSlice(slice0, tmp)
	if size > 11 {
		return slice1[0:12]
	}
	zeros := (12 - size)
	ztemp := []string{}
	for i := 0; i < zeros; i++ {
		ztemp = append(ztemp, "0")
	}
	return ConcatStringSlice(slice1, ztemp)
}

// CsvBucketWriter will append or create file to write out a WordBucket.
// general processing for golearn knn
func (wb *WordBucket) csvBucketKnnWriter(file_path string) {
	csvfile, err := os.OpenFile(file_path, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer csvfile.Close()

	writer := csv.NewWriter(csvfile)
	fit_value := fmt.Sprintf("%G", wb.FitValue) //0.003038961038961039
	bucketWrite := wb.BytePosCharToString()
	bucketWrite = append(bucketWrite, fit_value)
	bucketWrite = append(bucketWrite, wb.Word)
	bucketWrite = append(bucketWrite, wb.Category)
	writeErr := writer.Write(bucketWrite)

	/*
		//fit_value := fmt.Sprintf("%E", wb.FitValue) //3.038961E-03
		fit_value := fmt.Sprintf("%G", wb.FitValue) //0.003038961038961039
		bucket := fmt.Sprintf("%v", wb.Bucket)
		writeErr := writer.Write(
			[]string{
				bucket,
				fit_value,
				wb.Word,
				wb.Category,
			},
		)
	*/

	if writeErr != nil {
		fmt.Println(writeErr)
	}
	writer.Flush()
}

// CsvBucketWriter will append or create file to write out a WordBucket.
// Frozen limit of byters with headers to CSV, also includes fit_value.
// Process file in format for the knn_amitgupta_speedup.go script
func (wb *WordBucket) csvBucketKnnAmitWriter(file_path string) {
	csvfile, err := os.OpenFile(file_path, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer csvfile.Close()

	writer := csv.NewWriter(csvfile)

	//fit_value := fmt.Sprintf("%E", wb.FitValue) //3.038961E-03
	fit_value := fmt.Sprintf("%G", wb.FitValue) //0.003038961038961039
	bucketWrite := wb.BytePosCharToString()
	bucketWrite = append(bucketWrite, fit_value)
	writeErr := writer.Write(bucketWrite)

	if writeErr != nil {
		fmt.Println(writeErr)
	}
	writer.Flush()
}
