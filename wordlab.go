/*
* wordlab package pre-processes and formats words for further processing, specifically for classification or clustering algorithms (knn, k-means, x-means, etc...).
* It creates a unique floating point numeral for each unique word and writes to a file.
 */
package wordlab

type WordlabFormat interface {
	ParseInputWriteOut()
	WriteAttributes()
}
