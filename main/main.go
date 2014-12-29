package main

import (
	"fmt"
	"github.com/jbowles/wordlab"
	//"strings"
	tkz "github.com/jbowles/nlpt_tkz"
	"github.com/sjwhitworth/golearn/base"
	"strconv"
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
	HotelData()
}

func HotelData() {
	var root_errorfp = "/Users/jbowles/x/training_data/partner_fusion_trained_errors/training_data/"
	var root_datafp = "datasets/"
	wordlab.BuildHotelProviderDataKnn(root_errorfp, root_datafp)
	wordlab.BuildHotelProviderDataKnnAmit(root_errorfp, root_datafp)
}

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

/*
func makePrediction() base.FixedDataGrid {
	write_filep := "predict_dummy_file.csv"
	//write_filep := "datasets/wordlab_hotel_error_train.csv"
	headers := []string{"ByteRange0", "ByteRange1", "ByteRange2", "ByteRange3", "ByteRange4", "ByteRange5", "ByteRange6", "ByteRange7", "ByteRange8", "ByteRange9", "ByteRange10", "ByteRange11", "AggregateByteValue", "CategoryId"}
	wordlab.CsvCreateFileWithHeaders(true, write_filep, headers)

	sent := "unknown exception"
	tokens, _ := tkz.Tokenize("bukt", sent)
	for _, token := range tokens {
		bucket := wordlab.NewWordBucket(token, "")
		wordlab.CSVWriteAllAttributes(write_filep, bucket)
	}

	cls := wordlab.InitKnnClassifier(2, "euclidean", "datasets/wordlab_hotel_error_train.csv")
	fmt.Printf("cls: %v\n\n", cls)
	clsData, _ := base.ParseCSVToInstances(write_filep, true)
	fmt.Printf("clsData: %v\n\n", clsData)
	return cls.Predict(clsData)
	// -> Prediction: <nil>
}

func getBucket() *wordlab.WordBucket {
	return wordlab.NewWordBucket("reservation", "")
}
*/

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
