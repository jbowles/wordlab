package main

import (
	"fmt"
	"github.com/jbowles/wordlab"
	"runtime"
	"time"
	//"strings"
	//tkz "github.com/jbowles/nlpt_tkz"
	//ir "github.com/jbowles/nlpt_ir"
	"github.com/sjwhitworth/golearn/base"
	"github.com/sjwhitworth/golearn/evaluation"
	"github.com/sjwhitworth/golearn/knn"
	//"strconv"
)

func main() {
	//pipeFile()
	//pipeWildDir()
	//pipeHotelErrDir()
	//pipeWildDirOpt()
	//streamWildDirOpt()
	//streamHotelErrDirOpt()
	//HotelData()
	//wordlab.AmitKnnValidate()
	buildIt()
}

func streamHotelErrDirOpt() {
	runtime.GOMAXPROCS(8)
	//var space []ir.VecField
	tkzr := "lex"
	for docNum, value := range wordlab.HotelErrorIDTableDirs {
		go wordlab.StreamTokenizedDirectory(value[1], value[2], tkzr, docNum, time.Minute*90)
	}

	var input string
	fmt.Scanln(&input)
}

func buildIt() {
	wordlab.CsvCreateFileWithHeaders(true, "attributes.csv", []string{"term", "index", "dotproduct", "label"})
	for docNum, value := range wordlab.HotelErrorIDTableDirs {
		wordlab.BuildDocument(value[2], docNum)
	}
}

func pipeHotelErrDir() {
	//wordlab.CsvCreateFileWithHeaders(true, "attributes.csv", []string{"vetor", "index", "dotproduct", "label"})
	runtime.GOMAXPROCS(8)
	tkzr := "unicode"
	//tfidfFile := "datasets/tfidf.txt"
	//before goroutines this took == minutes to run against 2 directories of about 1000 files each.
	// goroutines this took == minutes to run against 2 directories of about 1000 files each.
	for docNum, value := range wordlab.HotelErrorIDTableDirs {
		wordlab.PipeTokenizedDirectory(value[1], value[2], tkzr, docNum, time.Minute*90)
	}
}

func streamWildDirOpt() {
	runtime.GOMAXPROCS(8)
	tkzr := "unicode"
	for docNum, value := range wordlab.News {
		go wordlab.StreamTokenizedDirectory(value[1], value[2], tkzr, docNum, time.Minute*90)
	}
	// sigkill goroutines in script
	var input string
	fmt.Scanln(&input)
}

//go run main/main.go  290.31s user 2.67s system 673% cpu 43.495 total
//go run main/main.go  145.14s user 1.09s system 693% cpu 21.080 total
func pipeWildDirOpt() {
	runtime.GOMAXPROCS(8)
	tkzr := "lex"
	for docNum, value := range wordlab.News {
		go wordlab.PipeTokenizedDirectoryOpt(value[1], value[2], tkzr, docNum, time.Minute*90)
	}
	var input string
	fmt.Scanln(&input)
}

func pipeWildDir() {
	//wordlab.CsvCreateFileWithHeaders(true, "attributes.csv", []string{"vetor", "index", "dotproduct", "label"})
	runtime.GOMAXPROCS(8)
	tkzr := "lex"
	//tfidfFile := "datasets/tfidf.txt"
	//before goroutines this took == minutes to run against 2 directories of about 1000 files each.
	// goroutines this took == minutes to run against 2 directories of about 1000 files each.
	for docNum, value := range wordlab.News {
		wordlab.PipeTokenizedDirectory(value[1], value[2], tkzr, docNum, time.Minute*90)
	}

	/*
		for _, value := range wordlab.HotelErrorIDTableDirs {
			wordlab.PipeTokenizedDirectory(value[1], value[2], tkzr, time.Second*90)
		}
	*/
}

func pipeFile() {
	path := "/Users/jbowles/x/training_data/corpora/20news-18828/alt.atheism/51060"
	tkzr := "lex"
	//tkzr := "unicode"
	res, err := wordlab.PipeTokenizedFile(path, tkzr)
	fmt.Println(string(res))
	fmt.Println(err)
}

func printHotelTable() {
	for id, s := range wordlab.HotelErrorIDTableFiles {
		fmt.Printf("\n HotelIds: %v, Errors: %v\n", id, s[0])
	}
}

func HotelData() {
	var root_datafp = "datasets/"
	wordlab.BuildHotelProviderDataKnnLabelIdFirst(root_datafp)
	wordlab.BuildHotelProviderDataKnnLabelNameLast(root_datafp)
}

// for word model... yck
/*
func amitClassify() {
	var stopList = wordlab.StopWords("datasets/stopwords/stopwords.txt")
	sent := "We’re sorry but we were unable to process your request. You may have temporarily exceeded your credit or debit card limit. Please choose a different card and try again"
	tokens, _ := tkz.TokenizeStr(sent, "unicode")
	buckets := []wordlab.WordBucket{}
	for _, token := range tokens {
		if !stopList.IsStopWord[token] {
			buckets = append(buckets, *wordlab.NewPredictionWordBucket(token))
		}
	}

	for _, bucket := range buckets {
		attrs := []float32{}
		for _, bpc := range bucket.Bucket {
			attrs = append(attrs, float32(bpc.ByteCharacter))
		}
		res := wordlab.AmitClassify(attrs)
		i, _ := strconv.Atoi(string(res))
		fmt.Printf("result for %v: %v %v\n", bucket.Word, string(res), wordlab.HotelErrorIDTableFiles[i][0])
	}
	//attrs := []float32{102, 97, 105, 108, 101, 100, 0, 0, 0, 0, 0, 0, 0.01652997688235604}
	//result := wordlab.AmitClassify(attrs)
	//fmt.Printf("result: %v\n", string(result))
	//fmt.Printf("result: %v\n", attrs)
}
*/

func makePrediction() {
	write_filep := "predict_dummy_file.csv"
	headers := wordlab.SentModelHeaders
	//headers = append(headers, "LabelId")
	headers = append(headers, "LabelName")
	wordlab.CsvCreateFileWithHeaders(true, write_filep, headers)
	var root_errorfp = "/Users/jbowles/x/training_data/partner_fusion_trained_errors/training_data/"

	var HotelErrorIDTable = map[int][]string{
		wordlab.AvailID: {"availability_error", "availability_error/availability_error.csv"},
		wordlab.BookID:  {"booking_error", "booking_error/booking_error.csv"},
	}
	for id, table := range HotelErrorIDTable {
		smodel := &wordlab.SentenceModel{
			InputFilePath:  root_errorfp + table[1],
			OutputFilePath: write_filep,
			LabelName:      table[0],
			Tokenizer:      "bukt",
			LabelID:        id,
			ForceOverwrite: true,
			LabelFirst:     false,
			LabelNameFirst: false,
			AddLabelName:   true,
			AddLabelID:     false,
		}
		smodel.ParseInputWriteOut()
	}

	cls := knn.NewKnnClassifier("euclidean", 1)
	rawData, perr0 := base.ParseCSVToInstances(write_filep, true)
	if perr0 != nil {
		panic(fmt.Sprintf("parse csv instances error: %s", perr0.Error()))
	}
	trainData, _ := base.InstancesTrainTestSplit(rawData, 0.60)
	cls.Fit(trainData)

	//fmt.Printf("cls: %v\n\n", cls)
	rawData, perr := base.ParseCSVToInstances(write_filep, true)
	//fmt.Printf("raw data: %v\n\n", rawData)
	if perr != nil {
		panic(fmt.Sprintf("parse csv instances error: %s", perr.Error()))
	}
	_, testData := base.InstancesTrainTestSplit(rawData, 0.60)

	predictions := cls.Predict(testData)
	fmt.Println(predictions)

	confusionMat, err := evaluation.GetConfusionMatrix(testData, predictions)
	if err != nil {
		panic(fmt.Sprintf("Unable to get confusion matrix: %s", err.Error()))
	}
	fmt.Println(evaluation.GetSummary(confusionMat))
}

func makeAttrs() {
	newInst := base.NewDenseInstances()
	fmt.Println("building new attributes... ")
	attrs := make([]base.Attribute, 2) //number of attributes
	attrs[0] = base.NewFloatAttribute("FitValue")
	attrs[1] = new(base.CategoricalAttribute)
	attrs[1].SetName("Category")

	fmt.Printf("defined attributes and categorical... %v\n\n", newInst)

	newSpecs := make([]base.AttributeSpec, len(attrs))
	for i, a := range attrs {
		fmt.Printf("i: %v\t a: %v\n\n\n", i, a)
		newSpecs[i] = newInst.AddAttribute(a)
	}
	fmt.Printf("added attributes... %v\n\n", newInst)

	newInst.Extend(100)

	fmt.Println(attrs[1].GetType())
	fmt.Println(attrs[1].GetName())
	fmt.Println(attrs[1].String())
	fmt.Println(attrs[0].GetType())
	fmt.Println(attrs[0].GetName())
	fmt.Println(attrs[0].String())
	//fmt.Println(attrs[0].GetSysValFromString("availability_error"))
}
