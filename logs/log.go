// This package is for logging runtime data.
package logs

import (
	"github.com/sirupsen/logrus"
	"os"
)
/*
"Logger" prints the log to the file specified
 */
func Logger(l Log) {
	f, err := os.OpenFile("logs/Logs/"+l.Filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		logrus.Fatal(err)
	}
	defer f.Close()
	logrus.SetLevel(logrus.TraceLevel)
	logrus.SetFormatter(l.Format)
	logrus.SetOutput(f)
	lg := logrus.WithFields(
		logrus.Fields{
			"FilePath":   l.FilePath,
			"LineNumber": l.LineNumber,
		})

	switch l.Level {
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
}
