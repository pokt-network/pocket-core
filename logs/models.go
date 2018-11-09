// This package is for logging runtime data.
package logs

import (
	"github.com/sirupsen/logrus"
)

type Log struct{ // TODO consider optimizations with enum like structure to give options to log creator
	Filename 	string 					`json:"filename"`
	Format 		logrus.Formatter 		`json:"format"`
	Level		logrus.Level 			`json:"Level"`
	// fields below
	FilePath	string 					`json:"filepath"`	// where the log message came from
	LineNumber 	string 					`json:"LineNumber"`	// specific line number from the message
	Message		string					`json:"message"`
}

