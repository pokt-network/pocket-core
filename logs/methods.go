// This package is for logging runtime data.
package logs

import (
	"github.com/pokt-network/pocket-core/config"
	"github.com/pokt-network/pocket-core/const"
	"github.com/pokt-network/pocket-core/util"
	"github.com/sirupsen/logrus"
	"os"
	"strconv"
	"time"
)

// "methods.go" describes custom logging functions.

/*
"NewLog" creates a custom log and calls the logger function.
 */
func NewLog(message string, level LogLevel, format LogFormat) {
	currentTime := time.Now()                            // get current time
	frame := util.Caller()                        		 // get the caller from util TODO: not cross platform windows issue
	if(frame==nil){
		panic("Frame from new log was nil")
	}
	log := &Log{}                                        // create a new log structure
	log.Name = currentTime.Format("2006.01.02 15:04:05") // set the current time in the specified format
	log.FunctionName = frame.Func.Name()                 // set the name of the function
	log.FilePath = frame.File                            // set thee path of the file
	log.Lev = level                                      // set the level
	log.Fmt = format                                     // set the format of the log (json)
	log.LineNumber = strconv.Itoa(frame.Line)            // set the line number
	log.Message = message                                // set the message
	Logger(*log)                                         // call the logger function to write to file
}

/*
"Logger" prints the log to data directory
 */
func Logger(l Log) {
	f, err := os.OpenFile(config.GetConfigInstance().Datadir+ // open/create the new log file
		_const.FILESEPARATOR+"logs"+_const.FILESEPARATOR+
		l.Name, os.O_WRONLY|os.O_CREATE|os.O_APPEND, os.ModePerm)
	if err != nil {										  // if there is an error then handle
		logrus.Fatal(err)
	}
	defer f.Close()										  // close the file once the function finishes
	logrus.SetFormatter(l.Fmt.format)					  // set the format from the log
	logrus.SetOutput(f)									  // set the output file
	lg := logrus.WithFields(							  // set custom fields
		logrus.Fields{
			"FilePath":   l.FilePath,					  // the path of the file the log was called from
			"LineNumber": l.LineNumber,					  // the line number of the file the log was called from
			"FunctionName": l.FunctionName,				  // the function name of the file the log was called
		})

	switch l.Lev.level {								  // switch statement for the logging levels
	case logrus.InfoLevel:								  // this is necessary to ensure the logrus dependency
		lg.Info(l.Message)								  // exists within the logs package only.
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
}
