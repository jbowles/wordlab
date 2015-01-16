package wordlab

import (
	"bufio"
	"bytes"
	"fmt"
	ir "github.com/jbowles/nlpt_ir"
	tkz "github.com/jbowles/nlpt_tkz"
	"gopkg.in/pipe.v2"
	"io"
	"os"
	"time"
)

// ReadFile reads data from the file at path and writes it to the
// pipe's stdout. I've hijacked the pipe projects ReadFile function
// and stuck a text tokenzer inside of it.
// The tokenizer used here MUST be 'lex' OR 'unicode'. The latter is the fastest but less flexible and comprehensive, while the former is not much slower it will return alot of symbols and punctuation. If all you need is "words" then use the 'unicode' tokenizer.
func PipeFileTokens(readFile, tokenizer string) pipe.Pipe {
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
		file, err := os.Open(s.Path(readFile))
		//defer file.Close()
		if err != nil {
			return err
		}
		scanner := bufio.NewScanner(file)
		bufferCache := new(bytes.Buffer)
		byteLining := []byte{'\n'} //newline padding bytes for writing to file
		//can lines
		for scanner.Scan() {
			bufferCache.Write(
				tkz.TokenizeBytes(scanner.Bytes(), tkzType).Bytes,
			)
			//follow each buffer write with a new line
			bufferCache.Write(byteLining)
		}

		//close file as soon as we can but no sooner.
		file.Close()
		//Log.Debug("streamBytes from tokenzier: %d", bufferCache.Len())
		_, err = io.Copy(s.Stdout, bufferCache)
		if err != nil {
			Log.Error("%s", err)
		}
		//file.Close()
		return err
	})
}

func StreamTokenizedDirectory(directoryPath, writeFile, tkzType string, docNum int, timeoutLimit time.Duration) {
	//overwrite the output file
	f, err := os.Create(writeFile)
	f.Close()
	if err != nil {
		Log.Error("%s", err)
		return
	}

	handler := NewFileHandler(directoryPath)
	go func(handler *FileHandler, timeoutLimit time.Duration, docNum int, tkzType, fileWrite string) {
		for _, file := range handler.FullFilePaths {
			p := pipe.Line(
				PipeFileTokens(file, tkzType),
				pipe.AppendFile(fileWrite, 0644),
			)
			_, err := pipe.CombinedOutputTimeout(p, timeoutLimit)
			if err != nil {
				Log.Error("pipe.CombinedOutputTimeout: %s %s", file, err)
			}
		}
	}(handler, timeoutLimit, docNum, tkzType, writeFile)
	Log.Notice("read %d files for directory %s", len(handler.FullFilePaths), handler.DirName)
}

func BuildDocument(file_path string, docNum int) *ir.VecField {
	file, err := os.Open(file_path)

	if err != nil {
		panic(err.Error())
	}

	defer file.Close()

	reader := bufio.NewReader(file)
	buffer := bytes.NewBuffer(make([]byte, 0))

	var chunk []byte
	var eol bool
	var str_array []string

	for {
		if chunk, eol, err = reader.ReadLine(); err != nil {
			break
		}
		buffer.Write(chunk)
		if !eol {
			str_array = append(str_array, buffer.String())
			buffer.Reset()
		}
	}
	Log.Error("length fo str array: %v", len(str_array))

	if err == io.EOF {
		err = nil
	}
	vf := &ir.VecField{}
	vf.Compose(str_array, docNum)
	WriteAttributes(*vf)
	return vf
}

func aggByteVal(term string) float64 {
	var seqTotal int
	for idx, rn := range term {
		seqTotal += idx + int(rn)
	}
	return float64(seqTotal)
}

func WriteAttributes(vf ir.VecField) {
	Log.Debug("************************************ writing attributes")
	csvfile, writer := fileWriterPipe("attributes.csv")
	defer csvfile.Close()

	for word, vectors := range vf.Space {
		//fmt.Printf("word: %v\n", word)
		//fmt.Printf("value: %v\n", value)
		for _, vector := range vectors {
			var bucketWrite []string
			//bucketWrite = append(bucketWrite, fmt.Sprintf("%v", vector))
			//bucketWrite = append(bucketWrite, fmt.Sprintf("%d", vector.BloomFilter))
			//bucketWrite = append(bucketWrite, fmt.Sprintf("%f", aggByteVal(word)))
			bucketWrite = append(bucketWrite, word)
			bucketWrite = append(bucketWrite, fmt.Sprintf("%d", vector.Index))
			bucketWrite = append(bucketWrite, fmt.Sprintf("%G", vector.DotProduct))
			//bucketWrite = append(bucketWrite, fmt.Sprintf("%d", vector.DocNum))
			bucketWrite = append(bucketWrite, HotelErrorIDTableDirs[vector.DocNum][0])
			writeErr := writer.Write(bucketWrite)
			if writeErr != nil {
				fmt.Println(writeErr)
			}
			writer.Flush()
		}
	}

}
