package main

import (
	"fmt"
	"github.com/jbowles/wordlab"
	//"strings"
	//tkz "github.com/jbowles/nlpt_tkz"
	"github.com/sjwhitworth/golearn/base"
	"github.com/sjwhitworth/golearn/evaluation"
	"github.com/sjwhitworth/golearn/knn"
	//"strconv"
)

func main() {

	/*
		sent := "Payment Failure happened. PAYMENT_FAILURE com.orbitz.tbs.model.txn.PaymentException: Bad Auth. Causing CreditAuthResult:Referral - call local BA number"
		res := wordlab.NewSentenceBucket(sent, "testing", "bukt", 45)
		totalsize := 0
		for _, bps := range res.Bucket {
			totalsize += len(bps.BytesSequence)
		}
		bpsts := res.BytePosSeqToString()
		fmt.Printf("byte pos seq to string: %v", bpsts)
		fmt.Printf("length: %v", len(bpsts))
		fmt.Printf("total length: %v", totalsize)
			fmt.Printf("SentnceBucket length %v\n", len(res.Bucket))
			for _, bps := range res.Bucket {
				fmt.Printf("Bucket length %v\n", len(bps.BytesSequence))
			}
			fmt.Printf("Sentence Bucket: %v\n", res.Bucket)
			fmt.Printf("sent bucket %v\n", res)
	*/
	//amitClassify()
	//HotelData()
	//makeAttrs()
	//newInstance()
	//makePrediction()
	for k, v := range wordlab.StopWords("default").IsStopWord {
		fmt.Printf("%v : %v\n", k, v)
	}
}

func printHotelTable() {
	for id, s := range wordlab.HotelErrorIDTable {
		fmt.Printf("\n HotelIds: %v, Errors: %v\n", id, s[0])
	}
}

func HotelData() {
	var root_errorfp = "/Users/jbowles/x/training_data/partner_fusion_trained_errors/training_data/"
	var root_datafp = "datasets/"
	//wordlab.BuildHotelProviderDataKnnLabelIdFirst(root_errorfp, root_datafp)
	wordlab.BuildHotelProviderDataKnnLabelNameLast(root_errorfp, root_datafp)
}

/*
func amitClassify() {
	var stopList = wordlab.StopWords("datasets/stopwords/stopwords.txt")
	sent := "We’re sorry but we were unable to process your request. You may have temporarily exceeded your credit or debit card limit. Please choose a different card and try again"
	tokens, _ := tkz.Tokenize("bukt", sent)
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
		fmt.Printf("result for %v: %v %v\n", bucket.Word, string(res), wordlab.HotelErrorIDTable[i][0])
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

/*
func getBucket() *wordlab.WordBucket {
	return wordlab.NewWordBucket("reservation", "")
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

func validateKnnClassifier() {
	neighbors := 12
	distance := "manhattan"
	training := "datasets/wordlab_knn_hotel_error_test.csv"
	knnClassifier := wordlab.InitKnnClassifier(neighbors, distance, training)

	//dense_copy := base.NewDenseCopy(testData)
	//fmt.Printf("copy of dense instance knnClasifier: %v\n", dense_copy)
	//sent := "We regret to inform you your credit card was declined."

	fmt.Printf("%v", knnClassifier)
	//fmt.Printf("%v \n\n\n", testData)
	//fmt.Printf("*** Prections:: %v\n\n", knnClassifier.Predict(testData))
}

func newInstance() {
	// There is a way of creating new Instances from scratch.
	// Inside an Instance, everything's stored as float64
	newData := make([]float64, 2)
	newData[0] = 1.0
	newData[1] = 0.0

	// Let's create some attributes
	attrs := make([]base.Attribute, 2)
	attrs[0] = base.NewFloatAttribute("Arbitrary Float Quantity")
	attrs[1] = new(base.CategoricalAttribute)
	attrs[1].SetName("Class")
	// Insert a standard class
	attrs[1].GetSysValFromString("A")

	// Now let's create the final instances set
	newInst := base.NewDenseInstances()

	// Add the attributes
	newSpecs := make([]base.AttributeSpec, len(attrs))
	for i, a := range attrs {
		newSpecs[i] = newInst.AddAttribute(a)
	}

	// Allocate space
	newInst.Extend(1)

	// Write the data
	newInst.Set(newSpecs[0], 0, newSpecs[0].GetAttribute().GetSysValFromString("1.0"))
	newInst.Set(newSpecs[1], 0, newSpecs[1].GetAttribute().GetSysValFromString("A"))

	fmt.Println(newInst)
}
*/
