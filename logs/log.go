// This package is for logging runtime data.
package logs

import (
	"github.com/pokt-network/pocket-core/config"
	"github.com/pokt-network/pocket-core/const"
	"github.com/sirupsen/logrus"
	"os"
)
/*
"Logger" prints the log to the file specified
 */
func Logger(l Log) {
	f, err := os.OpenFile(config.GetInstance().Datadir+_const.FILESEPARATOR+"logs"+_const.FILESEPARATOR+
		l.Name, os.O_WRONLY|os.O_CREATE|os.O_APPEND, os.ModePerm)
	if err != nil {
		logrus.Fatal(err)
	}
	defer f.Close()
	logrus.SetLevel(logrus.TraceLevel)
	logrus.SetFormatter(l.Fmt.format)
	logrus.SetOutput(f)
	lg := logrus.WithFields(
		logrus.Fields{
			"FilePath":   l.FilePath,
			"LineNumber": l.LineNumber,
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
}
