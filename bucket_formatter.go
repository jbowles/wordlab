package wordlab

import (
	"encoding/csv"
	"fmt"
	tkz "github.com/jbowles/nlpt-tkz"
	"io"
	"os"
	"strconv"
	"strings"
)

type WordModel struct {
	InputFilePath  string
	OutputFilePath string
	LabelName      string
	Tokenizer      string
	LabelID        int
	ForceOverwrite bool
	LabelFirst     bool
	LabelNameFirst bool
	AddLabelName   bool
	AddLabelID     bool
}

type SentenceModel struct {
	InputFilePath  string
	OutputFilePath string
	LabelName      string
	Tokenizer      string
	LabelID        int
	ForceOverwrite bool
	LabelFirst     bool
	LabelNameFirst bool
	AddLabelName   bool
	AddLabelID     bool
}

// SentenceModelLabelLast ParseInputWriteOut() does not need tokenizer as the tokenization is done at the time
// of creating NewSentenceBucket and computing the byte sequence ranges and aggregate byte values
func (sm SentenceModel) ParseInputWriteOut() {
	csvFile, err := os.Open(sm.InputFilePath)
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
			sm.WriteAttributes(
				NewSentenceBucket(
					row,
					sm.LabelName,
					sm.Tokenizer,
					sm.LabelID,
				),
			)
		}
	}
}

func (sm SentenceModel) WriteAttributes(sb *SentenceBucket) {
	csvfile, writer := fileWriter(sm.OutputFilePath)
	defer csvfile.Close()

	aggr_byte_value := fmt.Sprintf("%G", sb.AggregagteByteValue) //0.003038961038961039
	labelid := fmt.Sprintf("%d", sb.LabelID)                     // 3345
	//hashing := fmt.Sprintf("%d", sb.Hashing)                     // 4783264

	var bucketWrite []string
	if sm.LabelFirst {
		if sm.LabelNameFirst {
			bucketWrite = ConcatStringSlice([]string{sm.LabelName}, sb.BytePosSeqToString())
		}
		bucketWrite = ConcatStringSlice([]string{labelid}, sb.BytePosSeqToString())
	} else {
		bucketWrite = sb.BytePosSeqToString()
	}

	bucketWrite = append(bucketWrite, aggr_byte_value)
	//bucketWrite = append(bucketWrite, hashing)

	// IF need label to be in last position
	// Add label id and name at n-1 and n position
	if sm.AddLabelName && !sm.LabelFirst {
		// don't write labelid if prediction becomes too difficult
		//bucketWrite = append(bucketWrite, labelid)
		bucketWrite = append(bucketWrite, sm.LabelName)
	} else if sm.AddLabelID {
		bucketWrite = append(bucketWrite, labelid)
	}

	writeErr := writer.Write(bucketWrite)

	if writeErr != nil {
		fmt.Println(writeErr)
	}
	writer.Flush()
}

// tokenizer can be 'bukt', 'lex'
func (wm WordModel) ParseInputWriteOut() {
	csvFile, err := os.Open(wm.InputFilePath)
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
			tokens, _ := tkz.TokenizeStr(strings.ToLower(row), wm.Tokenizer)
			for _, token := range tokens {
				if !stopList.IsStopWord[token] {
					wm.WriteAttributes(NewWordBucket(token, wm.LabelName, wm.LabelID))
				}
			}
		}
	}
}

func (wm WordModel) WriteAttributes(wb *WordBucket) {
	csvfile, writer := fileWriter(wm.OutputFilePath)
	defer csvfile.Close()

	aggr_byte_value := fmt.Sprintf("%G", wb.AggregagteByteValue) //0.003038961038961039
	labelid := fmt.Sprintf("%d", wm.LabelID)                     // 3345

	var bucketWrite []string
	if wm.LabelFirst {
		if wm.LabelNameFirst {
			bucketWrite = ConcatStringSlice([]string{wm.LabelName}, wb.BytePosCharToString())
		}
		bucketWrite = ConcatStringSlice([]string{labelid}, wb.BytePosCharToString())
	} else {
		bucketWrite = wb.BytePosCharToString()
	}

	bucketWrite = append(bucketWrite, aggr_byte_value)

	// IF need label to be in last position
	// Add label id and name at n-1 and n position
	if wm.AddLabelName && !wm.LabelFirst {
		// don't write labelid if prediction becomes too difficult
		//bucketWrite = append(bucketWrite, labelid)
		bucketWrite = append(bucketWrite, wm.LabelName)
	} else if wm.AddLabelID {
		bucketWrite = append(bucketWrite, labelid)
	}

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
