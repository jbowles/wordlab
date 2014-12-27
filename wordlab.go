/*
* wordlab package pre-processes and formats words for further processing, specifically for classification or clustering algorithms (knn, k-means, x-means, etc...).
* It creates a unique floating point numeral for each unique word and writes to a file.
 */
package wordlab

import (
	"github.com/sjwhitworth/golearn/base"
	"github.com/sjwhitworth/golearn/knn"
	"strings"
)

const (
	Eqlzr = 0.13
)

var (
	posTotal float64
	chrTotal float64
)

type WordBucket struct {
	Bucket   []BytePosChar
	Word     string
	FitValue float64
	Category string
}

type BytePosChar struct {
	BytePosition  float64
	ByteCharacter float64
}

func NameFromFilePath(file_path string) string {
	res1 := strings.Split(file_path, "/")
	res2 := res1[len(res1)-1]
	return strings.Split(res2, ".")[0]
}

func setBytePosChar(position, character float64) BytePosChar {
	return BytePosChar{
		BytePosition:  position,
		ByteCharacter: character,
	}
}

func (wb *WordBucket) setFitValue() {
	for _, bpc := range wb.Bucket {
		chrTotal += bpc.ByteCharacter
		posTotal += bpc.BytePosition
	}
	wb.FitValue = ((posTotal * Eqlzr) / chrTotal) // * 100.0 // move 0.0033834892232614887 to 0.3383489223261489
}

func InitKnnClassifier(neighbors int, distance, data_file string) (*knn.KNNClassifier, base.FixedDataGrid) {
	//rawData, err := base.ParseCSVToInstances("../datasets/wordlab_knn_hotel_error_test.csv", true)
	rawData, err := base.ParseCSVToInstances(data_file, true)
	if err != nil {
		panic(err)
	}

	var cls = &knn.KNNClassifier{}

	switch distance {
	case "euclidean":
		cls = knn.NewKnnClassifier(distance, neighbors)
	case "manhattan":
		cls = knn.NewKnnClassifier(distance, neighbors)
	default:
		cls = knn.NewKnnClassifier("euclidean", 12)
	}

	_, testData := base.InstancesTrainTestSplit(rawData, 0.05)
	cls.Fit(rawData)

	return cls, testData
}
