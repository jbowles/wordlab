/*
* toying around with https://github.com/amitkgupta/nearest_neighbour/blob/master/golang-k-nn-speedup.go
 */
package wordlab

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"math"
	"runtime"
	"strconv"
)

type AmitKnnLabelWithFeatures struct {
	Label    []byte
	Features []float32
}

func NewAmitKnnLabelWithFeatures(parsedLine [][]byte) AmitKnnLabelWithFeatures {
	label := parsedLine[0]
	features := make([]float32, len(parsedLine)-1)

	for i, feature := range parsedLine {
		// skip label
		if i == 0 {
			continue
		}

		features[i-1] = byteSliceTofloat32(feature)
	}

	return AmitKnnLabelWithFeatures{label, features}
}

var newline = []byte("\n")
var comma = []byte(",")

func byteSliceTofloat32(b []byte) float32 {
	x, _ := strconv.ParseFloat(string(b), 32) //10, 8)
	return float32(x)
}

func parseCSVFile(filePath string) []AmitKnnLabelWithFeatures {
	fileContent, _ := ioutil.ReadFile(filePath)
	lines := bytes.Split(fileContent, newline)
	numRows := len(lines)

	labelsWithFeatures := make([]AmitKnnLabelWithFeatures, numRows-2)

	for i, line := range lines {
		// skip headers
		if i == 0 || i == numRows-1 {
			continue
		}

		labelsWithFeatures[i-1] = NewAmitKnnLabelWithFeatures(bytes.Split(line, comma))
	}

	return labelsWithFeatures
}

func squareDistanceWithBailout(features1, features2 []float32, bailout float32) (d float32) {
	//fmt.Printf("length %v, features1 %v \n", len(features1), features1)
	//fmt.Printf("length %v, features2, %v \n", len(features2), features2)
	for i := 0; i < len(features1); i++ {
		//fmt.Printf("*** i %v \n", i)
		//fmt.Printf("feature1 %v, feature2, %v \n", features1[0], features2[0])
		//fmt.Printf("feature1 %v, feature2, %v \n", features1[i], features2[i])
		x := features1[i] - features2[i]
		d += x * x

		if d > bailout {
			break
		}
	}

	return
}

func AmitClassify(features []float32) (label []byte) {
	//var trainingSample = parseCSVFile("datasets/trainingsample.csv")
	//var trainingSample = parseCSVFile("datasets/wordlab_amit_hotel_error_train.csv")
	var trainingSample = parseCSVFile("datasets/wordlab_hotel_error_sents_labelfirst_train.csv")
	label = trainingSample[0].Label
	d := squareDistanceWithBailout(features, trainingSample[0].Features, math.MaxFloat32)

	for _, row := range trainingSample {
		dNew := squareDistanceWithBailout(features, row.Features, d)

		if dNew < d {
			label = row.Label
			d = dNew
		}
	}

	return
}

func AmitKnnValidate() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	//validationSample := parseCSVFile("datasets/validationsample.csv")
	//var validationSample = parseCSVFile("datasets/wordlab_amit_hotel_error_validation.csv")
	var validationSample = parseCSVFile("datasets/validate.csv")

	var totalCorrect float32 = 0
	successChannel := make(chan float32)

	for _, test := range validationSample {
		go func(t AmitKnnLabelWithFeatures) {
			if string(t.Label) == string(AmitClassify(t.Features)) {
				successChannel <- 1
			} else {
				successChannel <- 0
			}
		}(test)
	}

	for i := 0; i < len(validationSample); i++ {
		totalCorrect += <-successChannel
	}

	fmt.Println(float32(totalCorrect) / float32(len(validationSample)))
}
