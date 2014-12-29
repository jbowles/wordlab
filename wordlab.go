/*
* wordlab package pre-processes and formats words for further processing, specifically for classification or clustering algorithms (knn, k-means, x-means, etc...).
* It creates a unique floating point numeral for each unique word and writes to a file.
 */
package wordlab

import (
	tkz "github.com/jbowles/nlpt_tkz"
	"github.com/sjwhitworth/golearn/base"
	"github.com/sjwhitworth/golearn/knn"
	"strconv"
	"strings"
)

// ByteRangeLimit defines what will be tolerated as a longest sequence of bytes per word.
// There is not need to capture all bytes of all words as the AggregateByteValue is computed
// over the whole range of bytes, but having a good order of bytes-per-word is good for classification.
// Remember: byte ranges are defined by slice indexes... will always be 1+ defined here
const (
	Eqlzr                   = 0.13
	ByteRangeWordModelLimit = 12
	ByteRangeSentModelLimit = 50
)

var (
	posTotal float64
	chrTotal float64
	seqTotal float64
	stopList = StopWords("/usr/local/mygo/src/github.com/jbowles/wordlab/datasets/stopwords/stopwords.txt")
)

type WordBucket struct {
	Bucket              []BytePosChar
	Word                string
	AggregagteByteValue float64
	LabelName           string
	LabelID             int
}

type SentenceBucket struct {
	Bucket              []BytesPosSeq
	Sentence            string
	AggregagteByteValue float64
	LabelName           string
	LabelID             int
}

type BytePosChar struct {
	BytePosition  float64
	ByteCharacter float64
}

type BytesPosSeq struct {
	BytesPosition float64
	BytesSequence []float64
}

func NewWordBucket(word, labelName string, labelId int) *WordBucket {
	bucket := &WordBucket{
		Word:      word,
		LabelName: labelName,
		LabelID:   labelId,
	}

	for pos, chr := range []byte(word) {
		bucket.Bucket = append(bucket.Bucket, setBytePosChar(float64(pos), float64(chr)))
	}
	bucket.setAggregateByteValue()
	return bucket
}

func NewPredictionWordBucket(word string) *WordBucket {
	bucket := &WordBucket{
		Word: word,
	}

	for pos, chr := range []byte(word) {
		bucket.Bucket = append(bucket.Bucket, setBytePosChar(float64(pos), float64(chr)))
	}
	bucket.setAggregateByteValue()
	return bucket
}

func NewSentenceBucket(sentence, labelName, tokenizer string, labelId int) *SentenceBucket {
	bucket := &SentenceBucket{
		Sentence:  sentence,
		LabelName: labelName,
		LabelID:   labelId,
	}

	tokens, _ := tkz.Tokenize(tokenizer, strings.ToLower(sentence))

	for pos, token := range tokens {
		if !stopList.IsStopWord[token] {
			var byteSeq = []float64{}
			for _, byt := range []byte(token) {
				byteSeq = append(byteSeq, float64(byt))
			}
			bucket.Bucket = append(bucket.Bucket, setBytesPosSeq(float64(pos), byteSeq))
		}
	}
	bucket.setAggregateByteValue()
	return bucket
}

func NewPredictionSentenceBucket(sentence, tokenizer string) *SentenceBucket {
	bucket := &SentenceBucket{
		Sentence: sentence,
	}

	tokens, _ := tkz.Tokenize(tokenizer, strings.ToLower(sentence))

	for pos, token := range tokens {
		if !stopList.IsStopWord[token] {
			var byteSeq = []float64{}
			for _, byt := range []byte(token) {
				byteSeq = append(byteSeq, float64(byt))
			}
			bucket.Bucket = append(bucket.Bucket, setBytesPosSeq(float64(pos), byteSeq))
		}
	}
	bucket.setAggregateByteValue()
	return bucket
}

// ConcatRuneSlice concatenates two slices of runes.
// Takes two arguments that must be slices of type String, order of args is not important.
// Returns a new slice which is the result of copying the 2 slices passed in a args.
//
//  slice_one := []rune("this, that, other")
//  slice_two := []rune("when, where, why")
//  ConcatRuneSlice(slice_one,slice_two)
func ConcatStringSlice(slice1, slice2 []string) []string {
	new_slice := make([]string, len(slice1)+len(slice2))
	copy(new_slice, slice1)
	copy(new_slice[len(slice1):], slice2)
	return new_slice
}

func ConcatFloat32Slice(slice1, slice2 []float32) []float32 {
	new_slice := make([]float32, len(slice1)+len(slice2))
	copy(new_slice, slice1)
	copy(new_slice[len(slice1):], slice2)
	return new_slice
}

// test that count is under limit, or if limit reached right number of zeros added over threshold.
func (wb *WordBucket) BytePosCharToString() []string {
	size := len(wb.Bucket)
	base := []string{}
	for _, i := range wb.Bucket {
		base = append(base, strconv.Itoa(int(i.ByteCharacter)))
	}
	if size > ByteRangeWordModelLimit {
		return base[0:ByteRangeWordModelLimit]
	}
	zeros := (ByteRangeWordModelLimit - size)
	ztemp := []string{}
	for i := 0; i < zeros; i++ {
		ztemp = append(ztemp, "0")
	}
	return ConcatStringSlice(
		base,
		makeZerosString(ByteRangeWordModelLimit, size),
	)
}

// test that count is under limit, or if limit reached right number of zeros added over threshold.
func (sb *SentenceBucket) BytePosSeqToString() (base []string) {
	size := 0
	tmp := []string{}
	for _, bps := range sb.Bucket {
		size += len(bps.BytesSequence)
		for _, i := range bps.BytesSequence {
			tmp = append(tmp, strconv.Itoa(int(i)))
		}
		base = tmp
	}
	if len(base) > ByteRangeSentModelLimit {
		base = base[0:ByteRangeSentModelLimit]
	}

	base = ConcatStringSlice(
		base,
		makeZerosString(ByteRangeSentModelLimit, size),
	)
	return
}

// test that count is under limit, or if limit reached right number of zeros added over threshold.
func (sb *SentenceBucket) BytePosSeqToFloat32() (base []float32) {
	size := 0
	tmp := []float32{}
	for _, bps := range sb.Bucket {
		size += len(bps.BytesSequence)
		for _, i := range bps.BytesSequence {
			tmp = append(tmp, float32(i))
		}
		base = tmp
	}
	if len(base) > ByteRangeSentModelLimit {
		base = base[0:ByteRangeSentModelLimit]
	}

	base = ConcatFloat32Slice(
		base,
		makeZerosFloat32(ByteRangeSentModelLimit, size),
	)
	return
}

func makeZerosString(limit, sliceSize int) (ztemp []string) {
	zeros := (limit - sliceSize)
	for i := 0; i < zeros; i++ {
		ztemp = append(ztemp, "0")
	}
	return
}

func makeZerosFloat32(limit, sliceSize int) (ztemp []float32) {
	zeros := (limit - sliceSize)
	for i := 0; i < zeros; i++ {
		ztemp = append(ztemp, float32(0))
	}
	return
}

func setBytePosChar(position, character float64) BytePosChar {
	return BytePosChar{
		BytePosition:  position,
		ByteCharacter: character,
	}
}

func setBytesPosSeq(position float64, sequence []float64) BytesPosSeq {
	return BytesPosSeq{
		BytesPosition: position,
		BytesSequence: sequence,
	}
}

func (wb *WordBucket) setAggregateByteValue() {
	for _, bpc := range wb.Bucket {
		chrTotal += bpc.ByteCharacter
		posTotal += bpc.BytePosition
	}
	wb.AggregagteByteValue = ((posTotal * Eqlzr) / chrTotal)
}

func (sb *SentenceBucket) setAggregateByteValue() {
	for _, bps := range sb.Bucket {
		for _, seq := range bps.BytesSequence {
			seqTotal += seq
		}
		posTotal += bps.BytesPosition
	}
	sb.AggregagteByteValue = ((posTotal * Eqlzr) / seqTotal)
}

func InitKnnClassifier(neighbors int, distance, data_file string) *knn.KNNClassifier {
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
		cls = knn.NewKnnClassifier("euclidean", 2)
	}

	cls.Fit(rawData)

	return cls
}
