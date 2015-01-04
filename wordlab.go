/*
* wordlab package pre-processes and formats words for further processing, specifically for classification or clustering algorithms (knn, k-means, x-means, etc...).
* It creates a unique floating point numeral for each unique word and writes to a file.
 */
package wordlab

import (
	"github.com/op/go-logging"
	"os"
)

type WordlabFormat interface {
	ParseInputWriteOut()
	WriteAttributes()
}

var Log = logging.MustGetLogger("example")
var LogFormat = logging.MustStringFormatter(
	"%{color}%{time:15:04:05.000} %{pid} %{shortfile} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}",
)

func init() {
	// For demo purposes, create two backend for os.Stderr.
	backend1 := logging.NewLogBackend(os.Stderr, "", 0)
	backend2 := logging.NewLogBackend(os.Stderr, "", 0)

	// For messages written to backend2 we want to add some additional
	// information to the output, including the used log level and the name of
	// the function.
	backend2Formatter := logging.NewBackendFormatter(backend2, LogFormat)

	// Only errors and more severe messages should be sent to backend1
	backend1Leveled := logging.AddModuleLevel(backend1)
	backend1Leveled.SetLevel(logging.ERROR, "")

	// Set the backends to be used.
	logging.SetBackend(backend1Leveled, backend2Formatter)
}

/*

	log.Debug("debug")
	log.Info("info")
	log.Notice("notice")
	log.Warning("warning")
	log.Error("err")
	log.Critical("crit")
*/
