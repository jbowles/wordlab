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

type CSVFormat interface {
	ParseInputWriteOut()
	WriteAttributes()
}

type WordModelLabelFirst struct {
	InputFilePath  string
	OutputFilePath string
	LabelName      string
	Tokenizer      string
	LabelID        int
	ForceOverwrite bool
}

type WordModelLabelLast struct {
	InputFilePath  string
	OutputFilePath string
	LabelName      string
	Tokenizer      string
	LabelID        int
	ForceOverwrite bool
}

type SentenceModelLabelFirst struct {
	InputFilePath  string
	OutputFilePath string
	LabelName      string
	Tokenizer      string
	LabelID        int
	ForceOverwrite bool
}

type SentenceModelLabelLast struct {
	InputFilePath  string
	OutputFilePath string
	LabelName      string
	Tokenizer      string
	LabelID        int
	ForceOverwrite bool
}

// SentenceModelLabelLast ParseInputWriteOut() does not need tokenizer as the tokenization is done at the time
// of creating NewSentenceBucket and computing the byte sequence ranges and aggregate byte values
func (smLast SentenceModelLabelLast) ParseInputWriteOut() {
	csvFile, err := os.Open(smLast.InputFilePath)
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
			smLast.WriteAttributes(NewSentenceBucket(row, smLast.LabelName, smLast.Tokenizer, smLast.LabelID))
		}
	}
}

// tokenizer can be 'bukt', 'lex'
func (smFirst SentenceModelLabelFirst) ParseInputWriteOut() {
	csvFile, err := os.Open(smFirst.InputFilePath)
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
			smFirst.WriteAttributes(NewSentenceBucket(row, smFirst.LabelName, smFirst.Tokenizer, smFirst.LabelID))
		}
	}
}

func (smLast SentenceModelLabelLast) WriteAttributes(sb *SentenceBucket) {
	csvfile, writer := fileWriter(smLast.OutputFilePath)
	defer csvfile.Close()

	// need label to be in last position
	aggr_byte_value := fmt.Sprintf("%G", sb.AggregagteByteValue) //0.003038961038961039
	label := fmt.Sprintf("%d", sb.LabelID)                       // 3345
	bucketWrite := sb.BytePosSeqToString()
	bucketWrite = append(bucketWrite, aggr_byte_value)
	bucketWrite = append(bucketWrite, label)
	writeErr := writer.Write(bucketWrite)

	if writeErr != nil {
		fmt.Println(writeErr)
	}
	writer.Flush()
}

func (smFirst SentenceModelLabelFirst) WriteAttributes(sb *SentenceBucket) {
	csvfile, writer := fileWriter(smFirst.OutputFilePath)
	defer csvfile.Close()

	label := fmt.Sprintf("%d", sb.LabelID)                  // 3345
	byte_value := fmt.Sprintf("%G", sb.AggregagteByteValue) //0.003038961038961039

	// need label to be in first position
	bucketWrite := ConcatStringSlice([]string{label}, sb.BytePosSeqToString())
	bucketWrite = append(bucketWrite, byte_value)
	writeErr := writer.Write(bucketWrite)

	if writeErr != nil {
		fmt.Println(writeErr)
	}
	writer.Flush()
}

// tokenizer can be 'bukt', 'lex'
func (wmLast WordModelLabelLast) ParseInputWriteOut() {
	csvFile, err := os.Open(wmLast.InputFilePath)
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
			tokens, _ := tkz.Tokenize(wmLast.Tokenizer, strings.ToLower(row))
			for _, token := range tokens {
				if !stopList.IsStopWord[token] {
					//bucket := NewWordBucket(token, wmLast.LabelName, category_id)
					wmLast.WriteAttributes(NewWordBucket(token, wmLast.LabelName, wmLast.LabelID))
				}
			}
		}
	}
}

// tokenizer can be 'bukt', 'lex'
func (wmFirst WordModelLabelFirst) ParseInputWriteOut() {
	csvFile, err := os.Open(wmFirst.InputFilePath)
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
			tokens, _ := tkz.Tokenize(wmFirst.Tokenizer, strings.ToLower(row))
			for _, token := range tokens {
				if !stopList.IsStopWord[token] {
					//bucket := NewWordBucket(token, wmFirst.LabelName, category_id)
					wmFirst.WriteAttributes(NewWordBucket(token, wmFirst.LabelName, wmFirst.LabelID))
				}
			}
		}
	}
}

func (wmLast WordModelLabelLast) WriteAttributes(wb *WordBucket) {
	csvfile, writer := fileWriter(wmLast.OutputFilePath)
	defer csvfile.Close()

	// need label to be in last position
	aggr_byte_value := fmt.Sprintf("%G", wb.AggregagteByteValue) //0.003038961038961039
	label := fmt.Sprintf("%d", wb.LabelID)                       // 3345
	bucketWrite := wb.BytePosCharToString()
	bucketWrite = append(bucketWrite, aggr_byte_value)
	bucketWrite = append(bucketWrite, label)
	writeErr := writer.Write(bucketWrite)

	if writeErr != nil {
		fmt.Println(writeErr)
	}
	writer.Flush()
}

func (wmFirst WordModelLabelFirst) WriteAttributes(wb *WordBucket) {
	csvfile, writer := fileWriter(wmFirst.OutputFilePath)
	defer csvfile.Close()

	label := fmt.Sprintf("%d", wb.LabelID)                  // 3345
	byte_value := fmt.Sprintf("%G", wb.AggregagteByteValue) //0.003038961038961039

	// need label to be in first position
	bucketWrite := ConcatStringSlice([]string{label}, wb.BytePosCharToString())
	bucketWrite = append(bucketWrite, byte_value)
	writeErr := writer.Write(bucketWrite)

	if writeErr != nil {
		fmt.Println(writeErr)
	}
	writer.Flush()
}

// fileWriter opens read/write appendable file OR creates one AND creates a csv file writer.
// NOTE: you must close the file 'defer osfile.Close()'.
func fileWriter(file_path string) (*os.File, *csv.Writer) {
	osfile, err := os.OpenFile(file_path, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println(err)
	}
	csvwriter := csv.NewWriter(osfile)
	return osfile, csvwriter
}

func CreateByteRangeHeaders(limit int) (hd []string) {
	base := "ByteRange"
	last := "AggregateByteValue"
	for i := 0; i < limit; i++ {
		hd = append(hd, (base + strconv.Itoa(i)))
	}
	hd = append(hd, last)
	return
}

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
