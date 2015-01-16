package wordlab

import (
	"bufio"
	"fmt"
	tkz "github.com/jbowles/nlpt_tkz"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"
)

// Highest level container for web content
type Collection struct {
	DocList       []*Document
	BuildTime     time.Duration
	Retrieved     int
	UnderSizedDoc int
}

// Core data container for web content,
//    including parsed content such as Terms and sentences
type Document struct {
	Id        int
	Label     string
	Tokens    []string
	Terms     []rune
	Lines     []string
	BuildTime time.Duration
}

func timeTracker(start time.Time, name string) {
	var elapsed = time.Since(start)
	fmt.Printf("For %q elapsed time = \t actual: %v\n", name, elapsed)
}

// Make Document
// Stores Scanner and parses text into words and sentences
// Also tracks time it took to build
// Returns Document channel
func makeNewDocument(id int, dChan chan Document, bscan *bufio.Scanner, tokenizer, label string) {
	t0 := time.Now()

	//fmt.Printf("length %v", len(textset))
	doc := Document{
		Id:    id,
		Label: label,
	}

	for bscan.Scan() {
		tokens, _ := tkz.TokenizeStr(bscan.Text(), tokenizer)
		for _, token := range tokens {
			doc.Tokens = append(doc.Tokens, token)
			for _, rn := range token {
				doc.Terms = append(doc.Terms, rn)
			}
		}
		doc.Lines = append(doc.Lines, bscan.Text())
	}
	doc.BuildTime = time.Since(t0) // shorthand for time.Now().Sub(t0)
	dChan <- doc
}

// Collection of documents from Web requests
// Builds NewDocument and adds it to Collection
// Tracks time it took to build
// Returns Collection
func (handler *FileHandler) MakeNewCollection() (coll Collection) {
	t0 := time.Now()
	doC := make(chan Document)
	count := 0
	small_file_count := 0

	// returns HttpResponse
	for _, f := range handler.FullFilePaths {
		bufScnr, file := ReadTextScanner(f)
		count += 1
		fi, _ := file.Stat()
		//smaller than 10 bytes
		if fi.Size() < 10 {
			small_file_count += 1
		}
		defer file.Close()

		go makeNewDocument(count, doC, bufScnr, handler.Tokenizer, handler.DocumentLabel)
		doc_reciever := <-doC
		coll.DocList = append(coll.DocList, &doc_reciever)
	}
	coll.Retrieved = count
	coll.UnderSizedDoc = small_file_count
	coll.BuildTime = time.Since(t0)
	return
}

func (handler *FileHandler) MakeNewMultiCollection() (coll Collection) {
	t0 := time.Now()
	doC := make(chan Document)
	count := 0
	small_file_count := 0

	// returns HttpResponse
	for _, f := range handler.FullFilePaths {
		bufScnr, file := ReadTextScanner(f)
		count += 1
		fi, _ := file.Stat()
		//smaller than 10 bytes
		if fi.Size() < 10 {
			small_file_count += 1
		}
		defer file.Close()

		go makeNewDocument(count, doC, bufScnr, handler.Tokenizer, strings.Split(fi.Name(), ".")[0])
		doc_reciever := <-doC
		coll.DocList = append(coll.DocList, &doc_reciever)
	}
	coll.Retrieved = count
	coll.UnderSizedDoc = small_file_count
	coll.BuildTime = time.Since(t0)
	return
}

func MakeNewCollections(handlers []FileHandler) []Collection {
	var collections []Collection
	for _, handler := range handlers {
		collections = append(collections, handler.MakeNewCollection())
	}
	return collections
}

func ReadTextString(path string) (s []string) {
	file, err := os.Open(path)
	if err != nil {
		log.Printf("Error opening file %v\n", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		s = append(s, scanner.Text())
	}
	return
}

func ReadTextScanner(path string) (*bufio.Scanner, *os.File) {
	file, err := os.Open(path)
	if err != nil {
		log.Printf("Error opening file %v\n", err)
	}
	//defer file.Close()
	return bufio.NewScanner(file), file
}

/*
ReadfileByte streams a text file and returns slice of bytes via ioutil

*/
func ReadFileByte(path string) []byte {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Printf("Error reading file %v\n", err)
	}
	return data
}

func MakeCollectionVis(coll *Collection) {
	size := len(coll.DocList)
	total_words := 0
	total_runes := 0
	total_sentences := 0
	total_labels := 0
	labelList := []string{}

	for _, doc := range coll.DocList {
		total_words += len(doc.Tokens)
		total_runes += len(doc.Terms)
		total_sentences += len(doc.Lines)
		labelList = append(labelList, doc.Label)
	}

	labelList = DeDupStrSlice(labelList)
	total_labels = len(labelList)
	total_unretrieved := size - coll.Retrieved
	avg_word := float64(total_runes) / float64(total_words)
	success_percent := float64(coll.Retrieved) / float64(size) * 100
	log.Printf(
		"\nCollection build time = %v \n Collection size (# of documents) = %d\n Total Labels = %d\n Labels = %v\n Total words = %d\n Total Runes(Unicode Characters) = %d\n Average Length of Words = %f\n Total Sentences = %d\n UnderSizedDocuments = %d\n Total Unretrieved = %d \n Total Retrieved = %d \n Success = %f percent \n \n",
		coll.BuildTime,
		size,
		total_labels,
		labelList,
		total_words,
		total_runes,
		avg_word,
		total_sentences,
		coll.UnderSizedDoc,
		total_unretrieved,
		coll.Retrieved,
		success_percent,
	)
}

func DeDupStrSlice(s []string) []string {
	result := []string{}
	seen := map[string]string{}
	for _, val := range s {
		if _, ok := seen[val]; !ok {
			result = append(result, val)
			seen[val] = val
		}
	}
	return result
}

// TFreqNoarm calculates term frequency for a specific token
func (doc *Document) TFreqNorm(tk string) float64 {
	var wc = float64(len(doc.Tokens))
	var counter float64
	for _, t := range doc.Tokens {
		switch t {
		case tk:
			counter += 1
		}
	}
	return counter / wc
}

// TFreq calculates term frequency for a specific token
func (doc *Document) TFreq(tk string) float64 {
	var timer time.Duration
	timer = time.Nanosecond
	tfreq := make(chan float64)
	var wc = float64(len(doc.Tokens))
	go func(doc *Document) {
		var counter float64
		for _, t := range doc.Tokens {
			switch t {
			case tk:
				counter += 1
			}
		}
		tfreq <- counter
	}(doc)

	for {
		select {
		case <-time.After(timer):
			//clog.Println(os.Stderr, "TFreq", log.Lshortfile)
			//log.Printf(" %v counting... ", timer)
		case res := <-tfreq:
			return res / wc
		}
	}
}

func (doc *Document) TypeFrequencyChan(tf chan<- []string) {
	this := []string{}
	var counter = 0
	for _, tok := range doc.Tokens {
		tok_freq := doc.TFreq(tok)
		counter += 1
		this = append(this, fmt.Sprintf("\n\nToken '%s', frequency %f, for Document: %s\n", tok, tok_freq, doc.Label))
	}
	tf <- this
}

func FrequencyPrinter(c <-chan []string) {
	for {
		fmt.Println(<-c)
	}
}

func PrintAllDocFreqs() {
	collection := MultiDocumentRunner(pf_err_dir, "lex")
	var freqChan chan []string = make(chan []string)
	for _, doc := range collection.DocList {
		go doc.TypeFrequencyChan(freqChan)
		go FrequencyPrinter(freqChan)
	}
	var input string
	fmt.Scanln(&input)
}

// multiple documents as files grouped under one directory with different labels
func MultiDocumentRunner(dirPath, tknz string) Collection {
	// 0.8 seconds for request and response timeouts
	fh := NewFileHandler(dirPath, tknz)
	return fh.MakeNewMultiCollection()
}

// multiple files grouped under one directory with same label
func DocumentRunner(dirLabel, dirPath, tknz string) Collection {
	// 0.8 seconds for request and response timeouts
	fh := NewDirHandler(dirPath, dirLabel, tknz)
	return fh.MakeNewCollection()
}

// multiple dirs with mutliple files, each dir is a label.
func CollectionRunner(labelPaths map[string]string, tknz string) []Collection {
	// 0.8 seconds for request and response timeouts
	var collections []Collection
	for dirLabel, dirPath := range labelPaths {
		fh := NewDirHandler(dirPath, dirLabel, tknz)
		collections = append(collections, fh.MakeNewCollection())
	}
	return collections
}

var golearn_dir = "/Users/jbowles/x/training_data/go-learn-project-data-sets/datasets"
var pf_err_dir = "/Users/jbowles/x/training_data/partner_fusion_trained_errors/training_data_txt"
var labelDirs = make(map[string]string)

func main() {
	//PrintAllDocFreqs()
	//labelDirs["athiesm"] = "/Users/jbowles/x/training_data/corpora/20news-18828/alt.atheism/"
	//labelDirs["cess"] = "/Users/jbowles/x/training_data/corpora/cess_esp/"
	//labelDirs["genesis"] = "/Users/jbowles/x/training_data/corpora/genesis"
	//labelDirs["genesis"] = "/Users/jbowles/x/training_data/corpora/genesis"

	//collection := DocumentRunner("athiesm", labelDirs["athiesm"], "lex")
	collection := MultiDocumentRunner(golearn_dir, "lex")
	//collection := MultiDocumentRunner(pf_err_dir)
	MakeCollectionVis(&collection)

	//colls := CollectionRunner(labelDirs)
	//for _, coll := range colls {
	//	MakeCollectionVis(&coll)
	//}

	//var readChan chan []string
	//for _, doc := range collection.DocList {
	//fmt.Printf("Id: %v, Label: %v, Tokens: %v, Terms: %v Lines: %v Build: %v\n", doc.Id, doc.Label, len(doc.Tokens), len(doc.Terms), len(doc.Lines), doc.BuildTime)
	//freq := doc.TFreq("error")
	//fmt.Printf("Term Frequency for 'error' %f in document: %s\n", freq, doc.Label)
	//}
}
