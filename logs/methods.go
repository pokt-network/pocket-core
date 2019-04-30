// This package is for logging runtime data.
package logs

import (
	"os"
	"runtime"
	"strconv"
	"time"

	"github.com/pokt-network/pocket-core/config"
	"github.com/pokt-network/pocket-core/const"
	"github.com/sirupsen/logrus"
)

// "Caller" returns the caller of the function that called it
func caller() *runtime.Frame {
	// make a slice of unsigned integers
	fpcs := make([]uintptr, 1)
	// skip 3 to get the original function call
	n := runtime.Callers(3, fpcs)
	if n == 0 {
		return nil
	}
	// get the information by calling callersframes on fpcs
	info, err := runtime.CallersFrames(fpcs).Next()
	if err {
		return nil
	}
	return &info
}

// "NewLog" creates a custom log and calls the logger function.
func NewLog(message string, level LogLevel, format LogFormat) error {
	currentTime := time.Now()
	// get the caller from util
	frame := caller()
	if frame == nil {
		panic("Frame from new log was nil")
	}
	// create a new log structure
	log := Log{}
	log.Name = currentTime.UTC().Format("2006-01-02T15-04-05")
	// set the name of the function
	log.FunctionName = frame.Func.Name()
	// set the path of the file
	log.FilePath = frame.File
	log.Lev = level
	// json
	log.Fmt = format
	// set the line number
	log.LineNumber = strconv.Itoa(frame.Line)
	log.Message = message
	if err := Logger(log); err != nil {
		return err
	}
	return nil
}

// "Logger" prints the log to data directory
func Logger(l Log) error {
	// open/create the new log file
	f, err := os.OpenFile(config.GlobalConfig().DD+_const.FILESEPARATOR+"logs"+_const.FILESEPARATOR+l.Name, os.O_WRONLY|os.O_CREATE|os.O_APPEND, os.ModePerm)
	if err != nil {
		return err
	}
	defer f.Close()
	logrus.SetFormatter(l.Fmt.format)
	logrus.SetOutput(f)
	lg := logrus.WithFields(
		logrus.Fields{
			"FilePath":     l.FilePath,     // the path of the file the log was called from
			"LineNumber":   l.LineNumber,   // the line number of the file the log was called from
			"FunctionName": l.FunctionName, // the function name of the file the log was called
		})

	switch l.Lev.level {
	case logrus.InfoLevel:
		lg.Info(l.Message)
	case logrus.DebugLevel:
		lg.Debug(l.Message)
	case logrus.FatalLevel:
		lg.Fatal(l.Message)
	case logrus.ErrorLevel:
		lg.Error(l.Message)
	case logrus.PanicLevel:
		lg.Panic(l.Message)
	case logrus.WarnLevel:
		lg.Warn(l.Message)
	case logrus.TraceLevel:
		lg.Trace(l.Message)
	default:
		lg.Info(l.Message)
	}
	return nil
}
