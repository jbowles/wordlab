package wordlab

import (
	"bufio"
	"bytes"
	//"io/ioutil"
	//"encoding/binary"
	//"encoding/csv"
	//"encoding/gob"
	//"fmt"
	//"fmt"
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
				//pipe.Filter(func(line []byte) bool { return stopList.IsStopWord[string(line)] }),
				pipe.AppendFile(fileWrite, 0644),
				//PipeFileTokens(fileWrite, "unicode"),
				//pipe.AppendFile(fileWrite, 0644),ReadDocBuildTfidf(fileWrite),
			)
			//_, err := pipe.CombinedOutputTimeout(p, timeoutLimit)
			output, err := pipe.CombinedOutputTimeout(p, timeoutLimit)
			if err != nil {
				Log.Error("pipe.CombinedOutputTimeout: %s %s", file, err)
			}
			vecField, err := ir.DecodeVectorStreamBytes(output)
			if err != nil {
				Log.Error("ir.DecodeVectorStreamBytes: %s", err)
			}
			WriteVector(vecField)

			/// *************** DEBUGGING ****************
			//Log.Debug("FILE: %v\n filter %v\n", file, string(output))
			//Log.Debug("FILE: %v\n tokens %v\n", file, string(output))
			//Log.Debug("VectorField %v", vecField)
			/// *************** DEBUGGING ****************
		}
	}(handler, timeoutLimit, docNum, tkzType, writeFile)
	Log.Notice("read %d files for directory %s", len(handler.FullFilePaths), handler.DirName)
}

func BuildStrArray(file_path string, docNum int) *ir.VecField {
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

	if err == io.EOF {
		err = nil
	}
	vf := &ir.VecField{}
	vf.Compose(str_array, docNum)
	return vf
}
