package unit

import (
	"flag"
	"os"
	"os/user"
	"testing"
	"time"

	"github.com/pokt-network/pocket-core/logs"
)

func TestJSONLogs(t *testing.T) {
	if err := logs.Log("Unit test for the json log functionality", logs.InfoLevel, logs.JSONLogFormat); err != nil {
		t.Fatalf(err.Error())
	}
}

func TestTextLogs(t *testing.T) {
	if err := logs.Log("Unit test for the text log functionality", logs.InfoLevel, logs.TextLogFormatter); err != nil {
		t.Fatalf(err.Error())
	}
}

func TestJSONLogsFileGeneration(t *testing.T) {

	if err := logs.Log("Unit test for the .json log file generation", logs.InfoLevel, logs.JSONLogFormat); err != nil {
		t.Fatalf(err.Error())
	}

	filepath := generateFilePath(t, ".json")

	if fileExists(filepath) {
		t.Logf("JSON log file %s exists", filepath)

	} else {
		t.Fatalf("Error JSON log file %s not being created", filepath)
	}

}

func TestTextLogsFileGeneration(t *testing.T) {

	// defines flag for logformat .log
	flag.Set("logformat", ".log")
	flag.Parse()

	if err := logs.Log("Unit test for the .log file generation", logs.InfoLevel, logs.TextLogFormatter); err != nil {
		t.Fatalf(err.Error())
	}

	filepath := generateFilePath(t, ".log")

	if fileExists(filepath) {
		t.Logf("Log file %s exists", filepath)

	} else {
		t.Fatalf("Error .log file %s not being created", filepath)
	}
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return true
}

func fileExistsOpen(filename string) bool {

	if _, err := os.OpenFile(filename, os.O_RDWR, 0644); err != nil {
		return false
	}
	return true
}

func generateFilePath(t *testing.T, prefix string) string {
	currentTime := time.Now()
	logName := currentTime.UTC().Format("2006-01-02T15-04-05")

	homeFolder, err := user.Current()
	if err != nil {
		t.Fatalf(err.Error())
	}

	filepath := homeFolder.HomeDir + "/.pocket/logs/" + logName + prefix
	return filepath
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
