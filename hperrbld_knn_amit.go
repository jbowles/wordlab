/*
* Hotel Provider Error Builder
* TEMPROARY: move it to the smallgear project.
 */
package wordlab

var knn_amit = "amit"

func BuildHotelProviderDataKnnAmit(root_errorfp, root_datafp string) {
	write_filep := root_datafp + "wordlab_amit_hotel_error_train.csv"
	headers := []string{"ByteRange0", "ByteRange1", "ByteRange2", "ByteRange3", "ByteRange4", "ByteRange5", "ByteRange6", "ByteRange7", "ByteRange8", "ByteRange9", "ByteRange10", "ByteRange11", "FitValue"}
	CsvCreateFileWithHeaders(true, write_filep, headers)

	avail_errorfp := (root_errorfp + "availability_error/availability_error.csv")
	book_errorfp := (root_errorfp + "booking_error/booking_error.csv")
	cancel_errorfp := (root_errorfp + "cancel_error/cancel_error.csv")
	cancelforb_errorfp := (root_errorfp + "cancel_forbidden_error/cancel_forbidden_error.csv")
	creditda_errorfp := (root_errorfp + "credit_data_error/credit_data_error.csv")
	creditde_errorfp := (root_errorfp + "credit_decline_error/credit_decline_error.csv")
	credits_errorfp := (root_errorfp + "credit_service_error/credit_service_error.csv")
	unexp_errorfp := (root_errorfp + "unexpected_response_error/unexpected_response_error.csv")

	avail := NameFromFilePath(avail_errorfp)
	book := NameFromFilePath(book_errorfp)
	cancel := NameFromFilePath(cancel_errorfp)
	cancelforb := NameFromFilePath(cancelforb_errorfp)
	creditda := NameFromFilePath(creditda_errorfp)
	creditde := NameFromFilePath(creditde_errorfp)
	credits := NameFromFilePath(credits_errorfp)
	unexp := NameFromFilePath(unexp_errorfp)

	ReadCsvFormatCsv(avail_errorfp, write_filep, avail, knn_amit)
	ReadCsvFormatCsv(book_errorfp, write_filep, book, knn_amit)
	ReadCsvFormatCsv(cancel_errorfp, write_filep, cancel, knn_amit)
	ReadCsvFormatCsv(cancelforb_errorfp, write_filep, cancelforb, knn_amit)
	ReadCsvFormatCsv(creditda_errorfp, write_filep, creditda, knn_amit)
	ReadCsvFormatCsv(creditde_errorfp, write_filep, creditde, knn_amit)
	ReadCsvFormatCsv(credits_errorfp, write_filep, credits, knn_amit)
	ReadCsvFormatCsv(unexp_errorfp, write_filep, unexp, knn_amit)
}
