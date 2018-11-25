package logs

import (
	"strconv"
	"time"

	"github.com/pokt-network/pocket-core/util"
)

func LogConstructorAndLog(message string, level LogLevel, format LogFormat) {
	currentTime := time.Now()
	f, t := util.Caller()
	filepath, ln := f.FileLine(t)
	log := &Log{}
	log.Name = currentTime.Format("2006.01.02 15:04:05")
	log.FunctionName = f.Name()
	log.FilePath = filepath
	log.Lev = level
	log.Fmt = format
	log.LineNumber = strconv.Itoa(ln)
	log.Message = message
	Logger(*log)
}
