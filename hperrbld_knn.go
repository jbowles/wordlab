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

var HotelErrorIDTable = map[int][]string{
	AvailID:         {"availability_error", "availability_error/availability_error.csv"},
	BookID:          {"booking_error", "booking_error/booking_error.csv"},
	CancelID:        {"cancel_error", "cancel_error/cancel_error.csv"},
	CancelForbidID:  {"cancel_forbidden_error", "cancel_forbidden_error/cancel_forbidden_error.csv"},
	CreditDataID:    {"credit_data_error", "credit_data_error/credit_data_error.csv"},
	creditDeclineID: {"credit_decline_error", "credit_decline_error/credit_decline_error.csv"},
	CreditServiceID: {"credit_service_error", "credit_service_error/credit_service_error.csv"},
	UnexpectID:      {"unexpected_response_error", "unexpected_response_error/unexpected_response_error.csv"},
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

var wordModelHeaders = CreateByteRangeHeaders(ByteRangeWordModelLimit)
var sentModelHeaders = CreateByteRangeHeaders(ByteRangeSentModelLimit)

func BuildHotelProviderDataKnn(root_errorfp, root_datafp string) {
	new_word_hdrs := append(wordModelHeaders, "Label")
	new_sent_hdrs := append(sentModelHeaders, "Label")
	CsvCreateFileWithHeaders(true, (root_datafp + "wordlab_hotel_error_words_labellast_train.csv"), new_word_hdrs)
	CsvCreateFileWithHeaders(true, (root_datafp + "wordlab_hotel_error_sents_labellast_train.csv"), new_sent_hdrs)

	for id, table := range HotelErrorIDTable {
		wmodel := &WordModelLabelLast{
			InputFilePath:  root_errorfp + table[1],
			OutputFilePath: root_datafp + "wordlab_hotel_error_words_labellast_train.csv",
			LabelName:      table[0],
			Tokenizer:      "bukt",
			LabelID:        id,
			ForceOverwrite: true,
		}
		smodel := &SentenceModelLabelLast{
			InputFilePath:  root_errorfp + table[1],
			OutputFilePath: root_datafp + "wordlab_hotel_error_sents_labellast_train.csv",
			LabelName:      table[0],
			Tokenizer:      "bukt",
			LabelID:        id,
			ForceOverwrite: true,
		}

		smodel.ParseInputWriteOut()
		wmodel.ParseInputWriteOut()
	}
}

func BuildHotelProviderDataKnnAmit(root_errorfp, root_datafp string) {
	new_word_hdrs := ConcatStringSlice([]string{"Label"}, wordModelHeaders)
	new_sent_hdrs := ConcatStringSlice([]string{"Label"}, sentModelHeaders)
	CsvCreateFileWithHeaders(true, (root_datafp + "wordlab_hotel_error_words_labelfirst_train.csv"), new_word_hdrs)
	CsvCreateFileWithHeaders(true, (root_datafp + "wordlab_hotel_error_sents_labelfirst_train.csv"), new_sent_hdrs)

	for id, table := range HotelErrorIDTable {
		wmodel := &WordModelLabelFirst{
			InputFilePath:  root_errorfp + table[1],
			OutputFilePath: root_datafp + "wordlab_hotel_error_words_labelfirst_train.csv",
			LabelName:      table[0],
			Tokenizer:      "bukt",
			LabelID:        id,
			ForceOverwrite: true,
		}
		smodel := &SentenceModelLabelFirst{
			InputFilePath:  root_errorfp + table[1],
			OutputFilePath: root_datafp + "wordlab_hotel_error_sents_labelfirst_train.csv",
			LabelName:      table[0],
			Tokenizer:      "bukt",
			LabelID:        id,
			ForceOverwrite: true,
		}

		smodel.ParseInputWriteOut()
		wmodel.ParseInputWriteOut()
	}
}
