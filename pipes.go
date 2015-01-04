package wordlab

/*
Pipes will stream data from files and tokenize in the stream.
*/
import (
	"bufio"
	"bytes"
	//"encoding/binary"
	"encoding/csv"
	//"encoding/gob"
	"fmt"
	ir "github.com/jbowles/nlpt_ir"
	tkz "github.com/jbowles/nlpt_tkz"
	"gopkg.in/pipe.v2"
	"io"
	"os"
	//"path/filepath"
	"time"
)

// ReadFile reads data from the file at path and writes it to the
// pipe's stdout. I've hijacked the pipe projects ReadFile function
// and stuck a text tokenzer inside of it.
// The tokenizer used here MUST be 'lex' OR 'unicode'. The latter is the fastest but less flexible and comprehensive, while the former is not much slower it will return alot of symbols and punctuation. If all you need is "words" then use the 'unicode' tokenizer.
func ReadFileAndTokenize(path, tokenizer string) pipe.Pipe {
	//so we don't fail becuase of bad tokenizer input
	var tkzType string
	switch tokenizer {
	case "unicode":
		tkzType = tokenizer
	default:
		tkzType = "lex"
	}
	//Log.Debug("Using tokenizer type: %s", tkzType)

	return pipe.TaskFunc(func(s *pipe.State) error {
		file, err := os.Open(s.Path(path))
		if err != nil {
			return err
		}
		scanner := bufio.NewScanner(file)
		bufferCache := new(bytes.Buffer)
		for scanner.Scan() {
			bufferCache.Write(
				tkz.TokenizeBytes(scanner.Bytes(), tkzType).Bytes,
			)
		}
		//Log.Debug("streamBytes from tokenzier: %d", bufferCache.Len())
		_, err = io.Copy(s.Stdout, bufferCache)
		file.Close()
		if err != nil {
			Log.Error("%s", err)
		}
		return err
	})
}

// PipeTokenizedFile streams data from a specified file, tokenizes text on the stream and returns []byte output and error. Error should return nil and []bytes should be greater than one.
// Since we are only dealing with one file the byte size returned should not be huge and so we simply return the content for the user to handle.
func PipeTokenizedFile(filePath, tkzType string) ([]byte, error) {
	Log.Debug("defining pipe.Line, prepare to stream ONE file: %s", filePath)
	p := pipe.Line(
		ReadFileAndTokenize(filePath, tkzType),
		//pipe.AppendFile("datasets/athiest.txt", 0644),
	)

	output, err := pipe.CombinedOutput(p)
	if err != nil {
		Log.Error("%s", err)
	}
	if len(output) < 20 {
		Log.Warning("Check filePath for: '%s' (use PipeTokenizedDirectory for directories). Check that file is not empty!!", filePath)
	}
	Log.Debug("pipe.Line streaming ONE file finished with byte size: %d", len(output))
	return output, err
}

// PipeTokenizedDirectory streams data from a specified file, tokenizes text on the stream and writes to intermediate file. It logs any errors encounterd.
// Since we are only dealing with a directory of n-number files it implements a timeout as well as not returning any content... instead it writes to output file.
// TODO: not finished....
func PipeTokenizedDirectory(directoryPath, fileWrite, tkzType string, timeoutLimit time.Duration) {

	f, err := os.Create(fileWrite)
	f.Close()
	if err != nil {
		Log.Error("%s", err)
		return
	}
	handler := NewFileHandler(directoryPath)
	for _, file := range handler.FullFilePaths {
		//Log.Debug("streaming file: %v", file)
		p := pipe.Line(
			ReadFileAndTokenize(file, tkzType),
			pipe.AppendFile(fileWrite, 0644),
			ReadDocBuildTfidf(fileWrite),
			//pipe.AppendFile("modelTFIDF", 0644),
		)
		output, err := pipe.CombinedOutputTimeout(p, timeoutLimit)
		if err != nil {
			Log.Error("pipe.CombinedOutputTimeout: %s %s", file, err)
		}
		vecField, err := ir.DecodeVectorStreamBytes(output)
		Log.Warning("DecodeVectorStream %v, %v\n", vecField, err)
	}
	Log.Notice("read %d files for directory %s", len(handler.FullFilePaths), handler.DirName)
}
func ReadDocBuildTfidf(path string) pipe.Pipe {
	return pipe.TaskFunc(func(s *pipe.State) error {
		file, err := os.Open(s.Path(path))
		if err != nil {
			return err
		}
		scanner := bufio.NewScanner(file)
		bufferCache := new(bytes.Buffer)
		vf := &ir.VecField{}
		for scanner.Scan() {
			vf.Compose([]string{scanner.Text()})
		}
		bufferCache.Write(vf.EncodeVectorStream(*bufferCache).ByteEncoding)
		/*
			Log.Info("%s", buffer.String())
			WriteAttributes(vf)
		*/
		_, err = io.Copy(s.Stdout, bufferCache)
		file.Close()
		if err != nil {
			Log.Error("%s", err)
		}
		return err
	})
}

/*
func BuildIndex(docs []string) {
	CsvCreateFileWithHeaders(true, "attributes.csv", []string{"vetor", "index", "dotproduct", "label"})
	Log.Debug("************************************ writing attributes")
	WriteAttributes(vf)
}
*/

func WriteAttributes(vf *ir.VecField) {
	csvfile, writer := fileWriterPipe("attributes.csv")
	defer csvfile.Close()

	for _, value := range vf.Space {
		//fmt.Printf("%v, count: %v\n", key, len(value))
		//fmt.Printf("%v\n", value)
		for _, vector := range value {
			var bucketWrite []string
			bucketWrite = append(bucketWrite, fmt.Sprintf("%v", vector))
			//bucketWrite = append(bucketWrite, fmt.Sprintf("%d", vector.BloomFilter))
			bucketWrite = append(bucketWrite, fmt.Sprintf("%d", vector.Index))
			bucketWrite = append(bucketWrite, fmt.Sprintf("%G", vector.DotProduct))
			//bucketWrite = append(bucketWrite, fmt.Sprintf("%d", vector.DocNum))
			bucketWrite = append(bucketWrite, News[vector.DocNum][0])
			writeErr := writer.Write(bucketWrite)
			if writeErr != nil {
				fmt.Println(writeErr)
			}
			writer.Flush()
		}
	}

}

func fileWriterPipe(file_path string) (*os.File, *csv.Writer) {
	osfile, err := os.OpenFile(file_path, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println(err)
	}
	csvwriter := csv.NewWriter(osfile)
	return osfile, csvwriter
}

var newsRootfp = "/Users/jbowles/x/training_data/corpora/20news-18828/"

const (
	AtheismID = iota + 1
	ComputerGraphicsID
	ComputerMsWindowsID
)

/*
	ComputerIbmPcHardwareID
	ComputerMacHardwareID
	MiscForSaleID
	RecAutosID
	RecMotorcyclesID
	RecSportBaseballID
	RecSportHockeyID
	SciCryptID
	SciElectronicsID
	SciMedID
	SciSpaceID
	SocReligionChristianID
	TalkPoliticsGunsID
	TalkPoliticsMideasetID
	TalkPoliticsMiscID
	TalkReligionMiscID
)
*/

var News = map[int][]string{
	AtheismID: {"athiesm", newsRootfp + "alt.atheism", "datasets/athiest.txt"},
	//ComputerGraphicsID:  {"graphics", newsRootfp + "comp.graphics", "datasets/graphics.txt"},
	//ComputerMsWindowsID: {"computermswindows", newsRootfp + "comp.os.ms-windows.misc", "datasets/computermswindows.txt"},
}

/*
	ComputerMsWindowsID:     {"", ""},
	ComputerIbmPcHardwareID: {"", ""},
	ComputerMacHardwareID:   {"", ""},
	MiscForSaleID:           {"", ""},
	RecAutosID:              {"", ""},
	RecMotorcyclesID:        {"", ""},
	RecSportBaseballID:      {"", ""},
	RecSportHockeyID:        {"", ""},
	SciCryptID:              {"", ""},
	SciElectronicsID:        {"", ""},
	SciMedID:                {"", ""},
	SciSpaceID:              {"", ""},
	SocReligionChristianID:  {"", ""},
	TalkPoliticsGunsID:      {"", ""},
	TalkPoliticsMideasetID:  {"", ""},
	TalkPoliticsMiscID:      {"", ""},
	TalkReligionMiscID:      {"", ""},
}
*/

func CsvCreateFileWithHeadersPipe(force bool, file_path string, headers []string) {
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
func fileExistPipe(file_path string) bool {
	_, err := os.OpenFile(file_path, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	return err == nil
}
