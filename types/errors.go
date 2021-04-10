package types

import (
	"bytes"
	"encoding/json"
	"fmt"
	"runtime"
	"strings"

	"github.com/pkg/errors"
	abci "github.com/tendermint/tendermint/abci/types"
)

// CodeType - ABCI code identifier within codespace
type CodeType uint32

// CodespaceType - codespace identifier
type CodespaceType string

// IsOK - is everything okay?
func (code CodeType) IsOK() bool {
	return code == CodeOK
}

// SDK error codes
const (
	// Base error codes
	CodeOK                    CodeType = 0
	CodeInternal              CodeType = 1
	CodeTxDecode              CodeType = 2
	CodeInvalidSequence       CodeType = 3
	CodeUnauthorized          CodeType = 4
	CodeInsufficientFunds     CodeType = 5
	CodeUnknownRequest        CodeType = 6
	CodeInvalidAddress        CodeType = 7
	CodeInvalidPubKey         CodeType = 8
	CodeUnknownAddress        CodeType = 9
	CodeInsufficientCoins     CodeType = 10
	CodeInvalidCoins          CodeType = 11
	CodeOutOfGas              CodeType = 12
	CodeMemoTooLarge          CodeType = 13
	CodeInsufficientFee       CodeType = 14
	CodeTooManySignatures     CodeType = 15
	CodeGasOverflow           CodeType = 16
	CodeNoSignatures          CodeType = 17
	CodeNegativeAmont         CodeType = 18
	CodeBurnStakedTokens      CodeType = 19
	CodeForceValidatorUnstake CodeType = 20
	CodeInvalidSlash          CodeType = 21
	CodeModuleAccountCreate   CodeType = 22
	CodeForbidden             CodeType = 23

	// CodespaceRoot is a codespace for error codes in this file only.
	// Notice that 0 is an "unset" codespace, which can be overridden with
	// Error.WithDefaultCodespace().
	CodespaceUndefined CodespaceType = ""
	CodespaceRoot      CodespaceType = "sdk"
)

func unknownCodeMsg(code CodeType) string {
	return fmt.Sprintf("unknown code %d", code)
}

// NOTE: Don't stringer this, we'll put better messages in later.
func CodeToDefaultMsg(code CodeType) string {
	switch code {
	case CodeInternal:
		return "internal error"
	case CodeTxDecode:
		return "tx parse error"
	case CodeInvalidSequence:
		return "invalid sequence"
	case CodeUnauthorized:
		return "unauthorized"
	case CodeInsufficientFunds:
		return "insufficient funds"
	case CodeUnknownRequest:
		return "unknown request"
	case CodeInvalidAddress:
		return "invalid address"
	case CodeInvalidPubKey:
		return "invalid pubkey"
	case CodeUnknownAddress:
		return "unknown address"
	case CodeInsufficientCoins:
		return "insufficient coins"
	case CodeInvalidCoins:
		return "invalid coins"
	case CodeOutOfGas:
		return "out of gas"
	case CodeMemoTooLarge:
		return "memo too large"
	case CodeInsufficientFee:
		return "insufficient fee"
	case CodeTooManySignatures:
		return "maximum numer of signatures exceeded"
	case CodeNoSignatures:
		return "no signatures supplied"
	default:
		return unknownCodeMsg(code)
	}
}

//--------------------------------------------------------------------------------
// All errors are created via constructors so as to enable us to hijack them
// and inject stack traces if we really want to.

// nolint
func ErrInternal(msg string) Error {
	return newErrorWithRootCodespace(CodeInternal, msg)
}
func ErrTxDecode(msg string) Error {
	return newErrorWithRootCodespace(CodeTxDecode, msg)
}
func ErrInvalidSequence(msg string) Error {
	return newErrorWithRootCodespace(CodeInvalidSequence, msg)
}
func ErrUnauthorized(msg string) Error {
	return newErrorWithRootCodespace(CodeUnauthorized, msg)
}
func ErrInsufficientFunds(msg string) Error {
	return newErrorWithRootCodespace(CodeInsufficientFunds, msg)
}
func ErrUnknownRequest(msg string) Error {
	return newErrorWithRootCodespace(CodeUnknownRequest, msg)
}
func ErrInvalidAddress(msg string) Error {
	return newErrorWithRootCodespace(CodeInvalidAddress, msg)
}
func ErrUnknownAddress(msg string) Error {
	return newErrorWithRootCodespace(CodeUnknownAddress, msg)
}
func ErrInvalidPubKey(msg string) Error {
	return newErrorWithRootCodespace(CodeInvalidPubKey, msg)
}
func ErrInsufficientCoins(msg string) Error {
	return newErrorWithRootCodespace(CodeInsufficientCoins, msg)
}
func ErrInvalidCoins(msg string) Error {
	return newErrorWithRootCodespace(CodeInvalidCoins, msg)
}
func ErrOutOfGas(msg string) Error {
	return newErrorWithRootCodespace(CodeOutOfGas, msg)
}
func ErrMemoTooLarge(msg string) Error {
	return newErrorWithRootCodespace(CodeMemoTooLarge, msg)
}
func ErrInsufficientFee(msg string) Error {
	return newErrorWithRootCodespace(CodeInsufficientFee, msg)
}
func ErrTooManySignatures(msg string) Error {
	return newErrorWithRootCodespace(CodeTooManySignatures, msg)
}
func ErrNoSignatures(msg string) Error {
	return newErrorWithRootCodespace(CodeNoSignatures, msg)
}
func ErrGasOverflow(msg string) Error {
	return newErrorWithRootCodespace(CodeGasOverflow, msg)
}
func ErrInvalidSlash(msg string) Error {
	return newErrorWithRootCodespace(CodeInvalidSlash, msg)
}
func ErrNegativeAmount(msg string) Error {
	return newErrorWithRootCodespace(CodeNegativeAmont, msg)
}
func ErrBurnStakedTokens(msg string) Error {
	return newErrorWithRootCodespace(CodeBurnStakedTokens, msg)
}
func ErrForceValidatorUnstake(msg string) Error {
	return newErrorWithRootCodespace(CodeForceValidatorUnstake, msg)
}
func ErrModuleAccountCreate(msg string) Error {
	return newErrorWithRootCodespace(CodeModuleAccountCreate, msg)
}
func ErrForbidden(msg string) Error {
	return newErrorWithRootCodespace(CodeForbidden, msg)
}

//----------------------------------------
// sdk Error type
type Error interface {
	Error() string
	Stacktrace() error
	Trace(offset int, format string, args ...interface{}) error
	Data() interface{}

	// convenience
	TraceSDK(format string, args ...interface{}) Error

	// set codespace
	WithDefaultCodespace(CodespaceType) Error

	Code() CodeType
	Codespace() CodespaceType
	ABCILog() string
	Result() Result
	QueryResult() abci.ResponseQuery
}

// NewError - create an error.
func NewError(codespace CodespaceType, code CodeType, format string, args ...interface{}) Error {
	return newError(codespace, code, format, args...)
}

func newErrorWithRootCodespace(code CodeType, format string, args ...interface{}) *sdkError {
	return newError(CodespaceRoot, code, format, args...)
}

func newError(codespace CodespaceType, code CodeType, format string, args ...interface{}) *sdkError {
	if format == "" {
		format = CodeToDefaultMsg(code)
	}
	return &sdkError{
		codespace: codespace,
		code:      code,
		data: FmtError{
			format: format,
			args:   args,
		},
	}
}

type sdkError struct {
	codespace  CodespaceType
	code       CodeType
	data       interface{}    // associated data
	stacktrace []uintptr      // first stack trace
	msgtraces  []msgtraceItem // all messages traced
}

// Implements Error.
func (err *sdkError) WithDefaultCodespace(cs CodespaceType) Error {
	codespace := err.codespace
	if codespace == CodespaceUndefined {
		codespace = cs
	}
	return &sdkError{
		codespace: cs,
		code:      err.code,
	}
}

// Captures a stacktrace if one was not already captured.
func (err *sdkError) Stacktrace() error {
	if err.stacktrace == nil {
		var offset = 3
		var depth = 32
		err.stacktrace = captureStacktrace(offset, depth)
	}
	return err
}

// Implements ABCIError.
// nolint: errcheck
func (err *sdkError) Trace(offset int, format string, args ...interface{}) error {
	msg := fmt.Sprintf(format, args...)
	return err.doTrace(msg, 0)
}

// Implements ABCIError.
// nolint: errcheck
func (err *sdkError) TraceSDK(format string, args ...interface{}) Error {
	msg := fmt.Sprintf(format, args...)
	return err.doTrace(msg, 0)
}

func (err *sdkError) doTrace(msg string, n int) *sdkError {
	pc, _, _, _ := runtime.Caller(n + 2) // +1 for doTrace().  +1 for the caller.
	// Include file & line number & msg.
	// Do not include the whole stack trace.
	err.msgtraces = append(err.msgtraces, msgtraceItem{
		pc:  pc,
		msg: msg,
	})
	return err
}

// Implements ABCIError.
func (err *sdkError) Error() string {
	return fmt.Sprintf(`ERROR:
Codespace: %s
Code: %d
Message: %#v
`, err.codespace, err.code, fmt.Sprintf("%v", err))
}

// Implements Error.
func (err *sdkError) Codespace() CodespaceType {
	return err.codespace
}

// Return the "data" of this error.
// Data could be used for error handling/switching,
// or for holding general error/debug information.
func (err *sdkError) Data() interface{} {
	return err.data
}

// Implements Error.
func (err *sdkError) Code() CodeType {
	return err.code
}

// Implements ABCIError.
func (err *sdkError) ABCILog() string {
	errMsg := err.Error()
	jsonErr := humanReadableError{
		Codespace: err.codespace,
		Code:      err.code,
		Message:   errMsg,
	}

	var buff bytes.Buffer
	enc := json.NewEncoder(&buff)
	enc.SetEscapeHTML(false)

	if err := enc.Encode(jsonErr); err != nil {
		panic(errors.Wrap(err, "failed to encode ABCI error log"))
	}

	return strings.TrimSpace(buff.String())
}

func (err *sdkError) Result() Result {
	return Result{
		Code:      err.Code(),
		Codespace: err.Codespace(),
		Log:       err.ABCILog(),
	}
}

// QueryResult allows us to return sdk.Error.QueryResult() in query responses
func (err *sdkError) QueryResult() abci.ResponseQuery {
	return abci.ResponseQuery{
		Code:      uint32(err.Code()),
		Codespace: string(err.Codespace()),
		Log:       err.ABCILog(),
	}
}

func (err *sdkError) Format(s fmt.State, verb rune) {
	switch verb {
	case 'p':
		_, _ = s.Write([]byte(fmt.Sprintf("%p", &err)))
	default:
		if s.Flag('#') {
			_, _ = s.Write([]byte("--= Error =--\n"))
			// Write data.
			_, _ = s.Write([]byte(fmt.Sprintf("Data: %#v\n", err.data)))
			// Write msg trace items.
			_, _ = s.Write([]byte(fmt.Sprintf("ProtoMsg Traces:\n")))
			for i, msgtrace := range err.msgtraces {
				_, _ = s.Write([]byte(fmt.Sprintf(" %4d  %s\n", i, msgtrace.String())))
			}
			// Write stack trace.
			if err.stacktrace != nil {
				_, _ = s.Write([]byte(fmt.Sprintf("Stack Trace:\n")))
				for i, pc := range err.stacktrace {
					fnc := runtime.FuncForPC(pc)
					file, line := fnc.FileLine(pc)
					_, _ = s.Write([]byte(fmt.Sprintf(" %4d  %s:%d\n", i, file, line)))
				}
			}
			_, _ = s.Write([]byte("--= /Error =--\n"))
		} else {
			// Write msg.
			_, _ = s.Write([]byte(fmt.Sprintf("%v", err.data)))
		}
	}
}

func (err *sdkError) ErrorMsg() string {
	return fmt.Sprintf(fmt.Sprintf("%v", err))
}

func (err sdkError) MarshalJSON() ([]byte, error) {
	errMsg := err.ErrorMsg()
	jsonErr := humanReadableError{
		Codespace: err.codespace,
		Code:      err.code,
		Message:   errMsg,
	}
	return json.Marshal(jsonErr)
}

//----------------------------------------
// REST error utilities

// appends a message to the head of the given error
func AppendMsgToErr(msg string, err string) string {
	msgIdx := strings.Index(err, "message\":\"")
	if msgIdx != -1 {
		errMsg := err[msgIdx+len("message\":\"") : len(err)-2]
		errMsg = fmt.Sprintf("%s; %s", msg, errMsg)
		return fmt.Sprintf("%s%s%s",
			err[:msgIdx+len("message\":\"")],
			errMsg,
			err[len(err)-2:],
		)
	}
	return fmt.Sprintf("%s; %s", msg, err)
}

// returns the index of the message in the ABCI Log
// nolint: deadcode unused
func mustGetMsgIndex(abciLog string) int {
	msgIdx := strings.Index(abciLog, "message\":\"")
	if msgIdx == -1 {
		panic(fmt.Errorf("invalid error format: %s", abciLog))
	}
	return msgIdx + len("message\":\"")
}

// parses the error into an object-like struct for exporting
type humanReadableError struct {
	Codespace CodespaceType `json:"codespace"`
	Code      CodeType      `json:"code"`
	Message   string        `json:"message"`
}

type msgtraceItem struct {
	pc  uintptr
	msg string
}

func (mti msgtraceItem) String() string {
	fnc := runtime.FuncForPC(mti.pc)
	file, line := fnc.FileLine(mti.pc)
	return fmt.Sprintf("%s:%d - %s",
		file, line,
		mti.msg,
	)
}

//----------------------------------------
// stacktrace & msgtraceItem

func captureStacktrace(offset int, depth int) []uintptr {
	var pcs = make([]uintptr, depth)
	n := runtime.Callers(offset, pcs)
	return pcs[0:n]
}

type FmtError struct {
	format string
	args   []interface{}
}

func (fe FmtError) Error() string {
	return fmt.Sprintf(fe.format, fe.args...)
}

func (fe FmtError) String() string {
	return fmt.Sprintf("FmtError{format:%v,args:%v}",
		fe.format, fe.args)
}

func (fe FmtError) Format() string {
	return fe.format
}
