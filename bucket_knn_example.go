/*
* Hotel Provider Error Builder
* TEMPROARY: move it to the smallgear project.
 */
package wordlab

// Labels should be distributed pretty far apart for knn algorithms to get accuracy.
// For example, in one implementation, shifting the iota (103,206,412,824,...,13184) instead of incrementing by 1 (1,2,3,4,...,8) boosted
// accuracy from low 70% to 99%
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

var hotelRootFp = "/Users/jbowles/x/training_data/partner_fusion_trained_errors/training_data/"
var HotelErrorIDTableDirs = map[int][]string{
	AvailID:         {"availability_error", hotelRootFp + "availability_error", "datasets/tmpavail.csv"},
	BookID:          {"booking_error", hotelRootFp + "booking_error", "datasets/tmpbook.csv"},
	CancelID:        {"cancel_error", hotelRootFp + "cancel_error", "datasets/tmpcancel.csv"},
	CancelForbidID:  {"cancel_forbidden_error", hotelRootFp + "cancel_forbidden_error", "datasets/tmpcnlfrbd.csv"},
	CreditDataID:    {"credit_data_error", hotelRootFp + "credit_data_error", "datasets/tmpcrdat.csv"},
	creditDeclineID: {"credit_decline_error", hotelRootFp + "credit_decline_error", "datasets/tmpcrdecl.csv"},
	CreditServiceID: {"credit_service_error", hotelRootFp + "credit_service_error", "datasets/tmpcserv.csv"},
	UnexpectID:      {"unexpected_response_error", hotelRootFp + "unexpected_response_error", "datasets/tmpunex.csv"},
}

var HotelErrorIDTableFiles = map[int][]string{
	AvailID:         {"availability_error", hotelRootFp + "availability_error/availability_error.csv", "datasets/tmpavail.csv"},
	BookID:          {"booking_error", hotelRootFp + "booking_error/booking_error.csv", "datasets/tmpbook.csv"},
	CancelID:        {"cancel_error", hotelRootFp + "cancel_error/cancel_error.csv", "datasets/tmpcancel.csv"},
	CancelForbidID:  {"cancel_forbidden_error", hotelRootFp + "cancel_forbidden_error/cancel_forbidden_error.csv", "datasets/tmpcnlfrbd.csv"},
	CreditDataID:    {"credit_data_error", hotelRootFp + "credit_data_error/credit_data_error.csv", "datasets/tmpcrdat.csv"},
	creditDeclineID: {"credit_decline_error", hotelRootFp + "credit_decline_error/credit_decline_error.csv", "datasets/tmpcrdecl.csv"},
	CreditServiceID: {"credit_service_error", hotelRootFp + "credit_service_error/credit_service_error.csv", "datasets/tmpcserv.csv"},
	UnexpectID:      {"unexpected_response_error", hotelRootFp + "unexpected_response_error/unexpected_response_error.csv", "datasets/tmpunex.csv"},
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

func BuildHotelProviderDataKnnLabelNameLast(root_errorfp, root_datafp string) {
	new_word_hdrs := append(wordModelHeaders, "labelname")

	new_sent_hdrs := append(SentModelHeaders, "lablename")

	CsvCreateFileWithHeaders(true, (root_datafp + "wordlab_bucket_hotel_error_words_labelnamelast_train.csv"), new_word_hdrs)
	CsvCreateFileWithHeaders(true, (root_datafp + "wordlab_bucket_hotel_error_sents_labelnamelast_train.csv"), new_sent_hdrs)

	for id, table := range HotelErrorIDTableFiles {
		// add word label name last
		wmodel := &WordModel{
			InputFilePath:  root_errorfp + table[1],
			OutputFilePath: root_datafp + "wordlab_bucket_hotel_error_words_labelnamelast_train.csv",
			LabelName:      table[0],
			Tokenizer:      "bukt",
			LabelID:        id,
			ForceOverwrite: true,
			LabelFirst:     false,
			LabelNameFirst: false,
			AddLabelName:   true,
			AddLabelID:     false,
		}

		// add sentence label name last
		smodel := &SentenceModel{
			InputFilePath:  root_errorfp + table[1],
			OutputFilePath: root_datafp + "wordlab_bucket_hotel_error_sents_labelnamelast_train.csv",
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
		wmodel.ParseInputWriteOut()
	}
}

func BuildHotelProviderDataKnnLabelIdFirst(root_errorfp, root_datafp string) {
	new_idword_hdrs := ConcatStringSlice([]string{"LabelId"}, wordModelHeaders)
	new_idsent_hdrs := ConcatStringSlice([]string{"LabelId"}, SentModelHeaders)
	CsvCreateFileWithHeaders(true, (root_datafp + "wordlab_hotel_error_words_labelidfirst_train.csv"), new_idword_hdrs)
	CsvCreateFileWithHeaders(true, (root_datafp + "wordlab_hotel_error_sents_labelidfirst_train.csv"), new_idsent_hdrs)

	for id, table := range HotelErrorIDTableFiles {
		wmodel := &WordModel{
			InputFilePath:  root_errorfp + table[1],
			OutputFilePath: root_datafp + "wordlab_hotel_error_words_labelidfirst_train.csv",
			LabelName:      table[0],
			Tokenizer:      "bukt",
			LabelID:        id,
			ForceOverwrite: true,
			LabelFirst:     true,
			AddLabelName:   false,
			LabelNameFirst: false,
		}
		smodel := &SentenceModel{
			InputFilePath:  root_errorfp + table[1],
			OutputFilePath: root_datafp + "wordlab_hotel_error_sents_labelidfirst_train.csv",
			LabelName:      table[0],
			Tokenizer:      "bukt",
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
