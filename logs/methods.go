package logs

import "github.com/sirupsen/logrus"

func LogConstructor(filename, filePath,
linenumber, message string, level logrus.Level, format logrus.Formatter) *Log{
	log:=&Log{}
	log.Filename=filename
	log.Level= level
	log.Format=format
	log.FilePath=filePath
	log.LineNumber=linenumber
	log.Message=message
	return log
}

func LogConstructorAndLog(filename, filePath,
linenumber, message string, level logrus.Level, format logrus.Formatter) {
	log:=&Log{}
	log.Filename=filename
	log.Level= level
	log.Format=format
	log.FilePath=filePath
	log.LineNumber=linenumber
	log.Message=message
	Logger(*log)
}
