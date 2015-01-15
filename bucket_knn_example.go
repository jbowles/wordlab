/*
* Hotel Provider Error Builder
* TEMPROARY: move it to the smallgear project.
 */
package wordlab

// Labels should be distributed pretty far apart for knn algorithms to get accuracy.
// For example, in one implementation, shifting the iota (103,206,412,824,...,13184) instead of incrementing by 1 (1,2,3,4,...,8) boosted
// accuracy from low 70% to 99%
/*
const (
	AvailID = 1 << iota * 103
	BookID
	CancelID
	CancelForbidID
	CreditDataID
	creditDeclineID
	CreditServiceID
	UnexpectID
)
*/

const (
	AvailID         = iota //= 20
	BookID                 //= 60
	CancelID               //= 100
	CancelForbidID         //= 140
	CreditDataID           //= 180
	creditDeclineID        //= 220
	CreditServiceID        //= 260
	UnexpectID             //= 300
)

var HotelRootFpCSV = "/Users/jbowles/x/training_data/partner_fusion_trained_errors/training_data_csv/"
var HotelRootFpTXT = "/Users/jbowles/x/training_data/partner_fusion_trained_errors/training_data_txt/"
var HotelErrorIDTableDirs = map[int][]string{
	AvailID:         {"availability", HotelRootFpTXT + "availability_error", "datasets/tmpavail.txt"},
	BookID:          {"booking", HotelRootFpTXT + "booking_error", "datasets/tmpbook.txt"},
	CancelID:        {"cancel", HotelRootFpTXT + "cancel_error", "datasets/tmpcancel.txt"},
	CancelForbidID:  {"cancelforbidden", HotelRootFpTXT + "cancel_forbidden_error", "datasets/tmpcnlfrbd.txt"},
	CreditDataID:    {"creditdata", HotelRootFpTXT + "credit_data_error", "datasets/tmpcrdat.txt"},
	creditDeclineID: {"creditdecline", HotelRootFpTXT + "credit_decline_error", "datasets/tmpcrdecl.txt"},
	CreditServiceID: {"creditservice", HotelRootFpTXT + "credit_service_error", "datasets/tmpcserv.txt"},
	UnexpectID:      {"unexpectedresponse", HotelRootFpTXT + "unexpected_response_error", "datasets/tmpunex.txt"},
}

var HotelErrorIDTableFiles = map[int][]string{
	AvailID:         {"availability_error", HotelRootFpCSV + "availability_error/availability_error.csv", "datasets/tmpavail.csv"},
	BookID:          {"booking_error", HotelRootFpCSV + "booking_error/booking_error.csv", "datasets/tmpbook.csv"},
	CancelID:        {"cancel_error", HotelRootFpCSV + "cancel_error/cancel_error.csv", "datasets/tmpcancel.csv"},
	CancelForbidID:  {"cancel_forbidden_error", HotelRootFpCSV + "cancel_forbidden_error/cancel_forbidden_error.csv", "datasets/tmpcnlfrbd.csv"},
	CreditDataID:    {"credit_data_error", HotelRootFpCSV + "credit_data_error/credit_data_error.csv", "datasets/tmpcrdat.csv"},
	creditDeclineID: {"credit_decline_error", HotelRootFpCSV + "credit_decline_error/credit_decline_error.csv", "datasets/tmpcrdecl.csv"},
	CreditServiceID: {"credit_service_error", HotelRootFpCSV + "credit_service_error/credit_service_error.csv", "datasets/tmpcserv.csv"},
	UnexpectID:      {"unexpected_response_error", HotelRootFpCSV + "unexpected_response_error/unexpected_response_error.csv", "datasets/tmpunex.csv"},
}
var HotelErrorNameTable = map[string]int{
	"availability_error":        AvailID,
	"booking_error":             BookID,
	"cancel_error":              CancelID,
	"cancel_forbidden_error":    CancelForbidID,
	"credit_data_error":         CreditDataID,
	"credit_decline_error":      creditDeclineID,
	"credit_service_error":      CreditServiceID,
	"unexpected_response_error": UnexpectID,
}

func AveragedLabelId() int {
	sum := int(0)
	for id, _ := range HotelErrorIDTableFiles {
		sum += id
	}
	return (sum / (len(HotelErrorIDTableFiles) * len(HotelErrorIDTableFiles)))
}

var wordModelHeaders = CreateByteRangeHeaders(ByteRangeWordModelLimit)
var SentModelHeaders = CreateByteRangeHeaders(ByteRangeSentModelLimit)
var wordsBucketNameLast = "wordlab_bucket_hotel_error_words_labelnamelast_train.csv"
var wordsBucketIdFirst = "wordlab_bucket_hotel_error_words_labelidfirst_train.csv"
var sentsBucketNameLast = "wordlab_bucket_hotel_error_sents_labelnamelast_train.csv"
var sentsBucketIdFirst = "wordlab_bucket_hotel_error_sents_labelidfirst_train.csv"

func BuildHotelProviderDataKnnLabelNameLast(root_datafp string) {
	new_word_hdrs := append(wordModelHeaders, "labelname")

	new_sent_hdrs := append(SentModelHeaders, "LableName")
	//new_sent_hdrs := append(SentModelHeaders, "Hashing")
	//new_sent_hdrs = append(new_sent_hdrs, "LableName")

	CsvCreateFileWithHeaders(true, (root_datafp + wordsBucketNameLast), new_word_hdrs)
	CsvCreateFileWithHeaders(true, (root_datafp + sentsBucketNameLast), new_sent_hdrs)

	for id, table := range HotelErrorIDTableFiles {
		// add word label name last
		wmodel := &WordModel{
			InputFilePath:  table[1],
			OutputFilePath: root_datafp + wordsBucketNameLast,
			LabelName:      table[0],
			Tokenizer:      "unicode",
			LabelID:        id,
			ForceOverwrite: true,
			LabelFirst:     false,
			LabelNameFirst: false,
			AddLabelName:   true,
			AddLabelID:     false,
		}

		// add sentence label name last
		smodel := &SentenceModel{
			InputFilePath:  table[1],
			OutputFilePath: root_datafp + sentsBucketNameLast,
			LabelName:      table[0],
			Tokenizer:      "unicode",
			LabelID:        id,
			ForceOverwrite: true,
			LabelFirst:     false,
			LabelNameFirst: false,
			AddLabelName:   true,
			AddLabelID:     false,
		}

		smodel.ParseInputWriteOut()
		wmodel.ParseInputWriteOut()
	}
}

func BuildHotelProviderDataKnnLabelIdFirst(root_datafp string) {
	new_idword_hdrs := ConcatStringSlice([]string{"LabelId"}, wordModelHeaders)
	new_idsent_hdrs := ConcatStringSlice([]string{"LabelId"}, SentModelHeaders)
	//new_idsent_hdrs = append(new_idsent_hdrs, "Hashing")

	CsvCreateFileWithHeaders(true, (root_datafp + wordsBucketIdFirst), new_idword_hdrs)
	CsvCreateFileWithHeaders(true, (root_datafp + sentsBucketIdFirst), new_idsent_hdrs)

	for id, table := range HotelErrorIDTableFiles {
		wmodel := &WordModel{
			InputFilePath:  table[1],
			OutputFilePath: root_datafp + wordsBucketIdFirst,
			LabelName:      table[0],
			Tokenizer:      "unicode",
			LabelID:        id,
			ForceOverwrite: true,
			LabelFirst:     true,
			AddLabelName:   false,
			LabelNameFirst: false,
		}
		smodel := &SentenceModel{
			InputFilePath:  table[1],
			OutputFilePath: root_datafp + sentsBucketIdFirst,
			LabelName:      table[0],
			Tokenizer:      "unicode",
			LabelID:        id,
			ForceOverwrite: true,
			LabelFirst:     true,
			AddLabelName:   false,
			LabelNameFirst: false,
		}

		smodel.ParseInputWriteOut()
		wmodel.ParseInputWriteOut()
	}
}
