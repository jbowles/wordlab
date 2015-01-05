package wordlab

import (
	tkz "github.com/jbowles/nlpt_tkz"
	"strconv"
	"strings"
)

// Buckets are an experiment to pre-process text into a formatter for various algorithms.
// A bucket is a specific type with Byte Position and Byte Character/Sequence and uses the byte value and its indexed position to create an aggregate byte value and stores the both the aggreate byt evalue and a bucket in a specific struct type.

// ByteRangeLimit defines what will be tolerated as a longest sequence of bytes per word.
// There is not need to capture all bytes of all words as the AggregateByteValue is computed
// over the whole range of bytes, but having a good order of bytes-per-word is good for classification.
// Remember: byte ranges are defined by slice indexes... will always be 1+ defined here
const (
	ByteRangeWordModelLimit = 12
	ByteRangeSentModelLimit = 20
)

type PositionTotal float64
type CharacterTotal float64
type SequenceTotal float64

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

	tokens, _ := tkz.TokenizeStr(strings.ToLower(sentence), tokenizer)

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

	tokens, _ := tkz.TokenizeStr(strings.ToLower(sentence), tokenizer)

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
	var chrTotal float64
	for _, bpc := range wb.Bucket {
		chrTotal += (bpc.ByteCharacter + bpc.BytePosition)
	}
	wb.AggregagteByteValue = chrTotal / 0.13
}

func (sb *SentenceBucket) setAggregateByteValue() {
	var seqTotal float64
	for _, bps := range sb.Bucket {
		for idx, seq := range bps.BytesSequence {
			seqTotal += (seq + float64(idx))
		}
	}
	sb.AggregagteByteValue = seqTotal / 0.13
}
