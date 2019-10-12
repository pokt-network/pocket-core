// This package is for logging runtime data.
package logs

import (
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strconv"
	"time"

	"github.com/pokt-network/pocket-core/config"
	"github.com/sirupsen/logrus"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
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

// "Log" creates a custom log and calls the logger function.
func Log(message string, level LogLevel, format LogFormat) error {
	currentTime := time.Now()
	// get the caller from util
	frame := caller()
	if frame == nil {
		panic("Frame from new log was nil")
	}
	// create a new log structure
	log := log{}
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
	if err := logger(log, format); err != nil {
		return err
	}
	return nil
}

// "logger" prints the log to data directory
func logger(l log, format LogFormat) error {
	filename := config.GlobalConfig().DD + string(filepath.Separator) + "logs" + string(filepath.Separator) + "pocket_core" + config.GlobalConfig().LogFormat

	f := &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    config.GlobalConfig().LogSize, // megabytes
		MaxBackups: config.GlobalConfig().LogBackups,
		MaxAge:     config.GlobalConfig().LogAge, //days
		Compress:   config.GlobalConfig().LogCompress}

	logrus.SetFormatter(l.Fmt.format)

	// We create empty log fields by default
	lg := logrus.WithFields(logrus.Fields{})

	// If using textformatter, we output to stdout
	if reflect.TypeOf(format.format) == reflect.TypeOf(&logrus.TextFormatter{}) {
		logrus.SetOutput(os.Stdout)

	// If we are using json, we just log to file .log or .json depending of logformat
	} else if reflect.TypeOf(format.format) == reflect.TypeOf(&logrus.JSONFormatter{}) {
		if config.GlobalConfig().LogFormat == ".json" {
			lg = logrus.WithFields(
				logrus.Fields{
					"FilePath":     l.FilePath,     // the path of the file the log was called from
					"LineNumber":   l.LineNumber,   // the line number of the file the log was called from
					"FunctionName": l.FunctionName, // the function name of the file the log was called
				})

			logrus.SetOutput(f)

		} else if config.GlobalConfig().LogFormat == ".log" {
			Textformatter := new(logrus.TextFormatter)
			Textformatter.TimestampFormat = "02-01-2006 15:04:05"
			Textformatter.FullTimestamp = true

			lg = logrus.WithFields(
				logrus.Fields{
					"FilePath":     l.FilePath,     // the path of the file the log was called from
					"LineNumber":   l.LineNumber,   // the line number of the file the log was called from
					"FunctionName": l.FunctionName, // the function name of the file the log was called
				})

			logrus.SetFormatter(Textformatter)

			logrus.SetOutput(f)
		}

	}
	return writeToLog(l, lg)

}

func writeToLog(l log, lg *logrus.Entry) error {
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
