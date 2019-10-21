package unit

import (
	"os"
	"testing"

	"github.com/pokt-network/pocket-core/config"
	"github.com/pokt-network/pocket-core/logs"
	utils "github.com/pokt-network/pocket-core/tests/unit/utils"
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
	args := []string{"./pocket_core", "--logformat", ".json"}
	utils.StartKillPocketCore(args, 15, "terminated", 500, t)

	if err := logs.Log("Unit test for the .json log file generation", logs.InfoLevel, logs.JSONLogFormat); err != nil {
		t.Fatalf(err.Error())
	}

	filepath := generateFilePath(t, ".json", "")

	if fileExists(filepath) {
		t.Logf("JSON log file %s exists", filepath)

	} else {
		t.Fatalf("Error JSON log file %s not being created", filepath)
	}

}

func TestJSONLogsFileGenerationCustomLogDir(t *testing.T) {
	args := []string{"./pocket_core", "--logformat", ".json", "--logdir", "./"}
	utils.StartKillPocketCore(args, 15, "terminated", 500, t)

	if err := logs.Log("Unit test for the .json log file generation", logs.InfoLevel, logs.JSONLogFormat); err != nil {
		t.Fatalf(err.Error())
	}

	filepath := generateFilePath(t, ".json", "./")

	if fileExists(filepath) {
		t.Logf("JSON log file %s exists", filepath)

	} else {
		t.Fatalf("Error JSON log file %s not being created", filepath)
	}

}

func TestTextLogsFileGeneration(t *testing.T) {
	args := []string{"./pocket_core", "--logformat", ".log"}
	utils.StartKillPocketCore(args, 15, "terminated", 500, t)

	if err := logs.Log("Unit test for the .log file generation", logs.InfoLevel, logs.TextLogFormatter); err != nil {
		t.Fatalf(err.Error())
	}

	filepath := generateFilePath(t, ".log", "")

	if fileExists(filepath) {
		t.Logf("Log file %s exists", filepath)

	} else {
		t.Fatalf("Error .log file %s not being created", filepath)
	}
}

func TestTextLogsFileGenerationCustomLogDir(t *testing.T) {
	args := []string{"./pocket_core", "--logformat", ".log", "--logdir", "./"}
	utils.StartKillPocketCore(args, 15, "terminated", 500, t)

	if err := logs.Log("Unit test for the .log file generation", logs.InfoLevel, logs.TextLogFormatter); err != nil {
		t.Fatalf(err.Error())
	}

	filepath := generateFilePath(t, ".log", "./")

	if fileExists(filepath) {
		t.Logf("Log file %s exists", filepath)

	} else {
		t.Fatalf("Error .log file %s not being created", filepath)
	}

	utils.DeleteTestBinary(t)
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

func generateFilePath(t *testing.T, prefix string, logdir string) string {
	logName := "pocket_core"

	filepath := logdir + logName + prefix

	// If no logdir given, we assume logdir is on our default datadir of pocket config
	if len(logdir) == 0 {
		homedir := config.GlobalConfig().LogDir
		filepath = homedir + logName + prefix

	}

	return filepath
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
