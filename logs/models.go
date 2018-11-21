// This package is for logging runtime data.
package logs

import (
	"github.com/logmatic/logmatic-go"
	"github.com/sirupsen/logrus"
)

// "models.go" describes custom logging models

var (
	InfoLevel     = LogLevel{logrus.InfoLevel}					// wrappers for each level
	DebugLevel    = LogLevel{logrus.DebugLevel}				// this is necessary to ensure
	WaringLevel   = LogLevel{logrus.WarnLevel}					// the logrus dependency
	PanicLevel    = LogLevel{logrus.PanicLevel}				// exists only within the logs
	FatalLevel    = LogLevel{logrus.FatalLevel}				// package.
	ErrorLevel    = LogLevel{logrus.ErrorLevel}
	TraceLevel    = LogLevel{logrus.TraceLevel}
	JSONLogFormat = LogFormat{&logmatic.JSONFormatter{}} 	// wrapper for the json log format
)

/*
"Log" model holds the structure for the log properties.
 */
type Log struct {
	Name         string    `json:"filename"`						// name of the log file
	Fmt          LogFormat `json:"format"`							// format of the log
	Lev          LogLevel  `json:"Lev"`								// level of the log (see var above)
	FilePath     string    `json:"filepath"`   						// where the log message came from
	FunctionName string    `json:"functionname"`					// the functionName where the
	LineNumber   string    `json:"LineNumber"` 						// specific line number from the message
	Message      string    `json:"message"`							// the main message "payload" of the log.
}
/*
"LogLevel" model is a simple wrapper structure for the logrus level.
 */
type LogLevel struct {
	level logrus.Level
}
/*
"LogFormat" model is a simple wrapper structure for the logrus format.
 */
type LogFormat struct {
	format logrus.Formatter
}
