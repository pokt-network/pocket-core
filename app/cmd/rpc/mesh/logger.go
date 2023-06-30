package mesh

import (
	"fmt"
	kitlevel "github.com/go-kit/kit/log/level"
	"github.com/go-kit/kit/log/term"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/pokt-network/pocket-core/app"
	pocketTypes "github.com/pokt-network/pocket-core/x/pocketcore/types"
	"github.com/tendermint/tendermint/libs/cli/flags"
	"github.com/tendermint/tendermint/libs/log"
	log2 "log"
	"net/url"
	"os"
)

type LevelHTTPLogger struct {
	retryablehttp.LeveledLogger
}

var (
	LogLvlInfo  = "info"
	LogLvlError = "error"
	LogLvlDebug = "debug"
)

// fields - mutate interface to key/value object to be print on stdout
func (l *LevelHTTPLogger) fields(keysAndValues ...interface{}) map[string]interface{} {
	fields := make(map[string]interface{})

	for i := 0; i < len(keysAndValues)-1; i += 2 {
		fields[keysAndValues[i].(string)] = keysAndValues[i+1]
	}

	return fields
}

// Error - log to stdout as error level
func (l *LevelHTTPLogger) Error(msg string, keysAndValues ...interface{}) {
	fields := l.fields(keysAndValues...)
	err := fields["error"].(error)
	_url := fields["url"]
	if _url != nil {
		_url2, ok := _url.(*url.URL)
		if !ok {
			logger.Error("request error", "error", _url)
			return
		}

		logger.Error(
			fmt.Sprintf(
				"%s at %s %s://%s%s\n",
				msg,
				fields["method"].(string),
				_url2.Scheme,
				_url2.Host,
				_url2.Path,
			),
		)
		return
	}
	logger.Error(msg, err, fields)
}

// Info - log to stdout as info level
func (l *LevelHTTPLogger) Info(msg string, keysAndValues ...interface{}) {
	logger.Info(msg, l.fields(keysAndValues...))
}

// Debug - log to stdout as debug level
func (l *LevelHTTPLogger) Debug(msg string, keysAndValues ...interface{}) {
	fields := l.fields(keysAndValues...)
	_url := fields["url"]
	if _url != nil {
		_url2, ok := _url.(*url.URL)
		if !ok {
			logger.Error(fmt.Sprintf("unable to cast to url.URL %v", _url))
			return
		}
		logger.Debug(
			fmt.Sprintf(
				"%s:\nURL=%s://%s%s?%s\nMETHOD=%s",
				msg,
				_url2.Scheme, _url2.Host, _url2.Path, _url2.RawQuery,
				fields["method"].(string),
			),
		)
		return
	}
	logger.Debug(msg, fields)
}

// Warn - log to stdout as warning level
func (l *LevelHTTPLogger) Warn(msg string, keysAndValues ...interface{}) {
	logger.Debug(msg, l.fields(keysAndValues...))
}

// Log - log base on level
func Log(msg string, level string) {
	switch level {
	case LogLvlInfo:
		logger.Info(msg)
	case LogLvlDebug:
		logger.Debug(msg)
	case LogLvlError:
		logger.Error(msg)
	default:
		logger.Info(msg)
	}
}

func relayToString(relay *pocketTypes.Relay) string {
	sessionHash := sessionStorage.GetSessionHashFromRelay(relay)
	servicerAddress, _ := GetAddressFromPubKeyAsString(relay.Proof.ServicerPubKey)
	return fmt.Sprintf(
		"session_hash=%s session_height=%d app=%s chain=%s servicer=%s",
		sessionHash, relay.Proof.SessionBlockHeight,
		relay.Proof.Token.ApplicationPublicKey, relay.Proof.Blockchain,
		servicerAddress,
	)
}

func LogRelay(relay *pocketTypes.Relay, msg, level string) {
	Log(fmt.Sprintf("msg=%s %s", msg, relayToString(relay)), level)
}

// initLogger - initialize logger
func initLogger() (logger log.Logger) {
	logger = log.NewTMLoggerWithColorFn(log.NewSyncWriter(os.Stdout), func(keyvals ...interface{}) term.FgBgColor {
		if keyvals[0] != kitlevel.Key() {
			fmt.Printf("expected level key to be first, got %v", keyvals[0])
			log2.Fatal(1)
		}
		switch keyvals[1].(kitlevel.Value).String() {
		case LogLvlInfo:
			return term.FgBgColor{Fg: term.Green}
		case LogLvlDebug:
			return term.FgBgColor{Fg: term.DarkBlue}
		case LogLvlError:
			return term.FgBgColor{Fg: term.Red}
		default:
			return term.FgBgColor{Fg: term.Green}
		}
	})
	l, err := flags.ParseLogLevel(app.GlobalMeshConfig.LogLevel, logger, "info")
	if err != nil {
		log2.Fatal(err)
	}
	logger = l
	return
}

func GetLogger() log.Logger {
	return logger
}
