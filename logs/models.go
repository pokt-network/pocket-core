// This package is for logging runtime data.
package logs

import (
	"github.com/logmatic/logmatic-go"
	"github.com/sirupsen/logrus"
)

var (
	InfoLevel     = LogLevel{logrus.InfoLevel}
	DebugLevel    = LogLevel{logrus.DebugLevel}
	WaringLevel   = LogLevel{logrus.WarnLevel}
	PanicLevel    = LogLevel{logrus.PanicLevel}
	FatalLevel    = LogLevel{logrus.FatalLevel}
	ErrorLevel    = LogLevel{logrus.ErrorLevel}
	TraceLevel    = LogLevel{logrus.TraceLevel}
	JSONLogFormat = LogFormat{&logmatic.JSONFormatter{}}
)

type Log struct {
	Name         string    `json:"filename"`
	Fmt          LogFormat `json:"format"`
	Lev          LogLevel  `json:"Lev"`
	FilePath     string    `json:"filepath"`   // where the log message came from
	FunctionName string    `json:"functionname"`
	LineNumber   string    `json:"LineNumber"` // specific line number from the message
	Message      string    `json:"message"`
}

// TODO consider optimizations with enum like structure to give options to log creator

type LogLevel struct {
	level logrus.Level
}

type LogFormat struct {
	format logrus.Formatter
}
