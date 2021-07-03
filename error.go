package errors

import (
	"fmt"
	"github.com/valyala/bytebufferpool"
	"runtime"
)

const (
	invalidArgumentErrorFailureCodeCode = 400
	invalidArgumentErrorCode            = "***BAD REQUEST***"

	unauthorizedErrorFailureCodeCode = 401
	unauthorizedErrorCode            = "***UNAUTHORIZED***"

	forbiddenErrorFailureCodeCode = 403
	forbiddenErrorCode            = "***FORBIDDEN***"

	notFoundErrorFailureCodeCode = 404
	notFoundErrorCode            = "***NOT FOUND***"

	serviceErrorFailureCodeCode = 500
	serviceErrorCode            = "***SERVICE EXECUTE FAILED***"

	serviceNotImplementedErrorFailureCodeCode = 501
	serviceNotImplementedErrorCode            = "***SERVICE NOT IMPLEMENTED***"

	unavailableErrorFailureCodeCode = 503
	unavailableErrorCode            = "***SERVICE UNAVAILABLE***"
)

type CodeError interface {
	SetId(id string) CodeError
	SetFailureCode(failureCode int) CodeError
	GetMeta() MultiMap
	GetStacktrace() (fn string, file string, line int)
	Error() string
	String() string
	ToJson() []byte
}

type stacktrace struct {
	Fn   string `json:"fn"`
	File string `json:"file"`
	Line int    `json:"line"`
}

type codeError struct {
	Id          string     `json:"id,omitempty"`
	FailureCode int        `json:"failureCode,omitempty"`
	Code        string     `json:"code,omitempty"`
	Message     string     `json:"message,omitempty"`
	Meta        MultiMap   `json:"meta,omitempty"`
	Stacktrace  stacktrace `json:"stacktrace,omitempty"`
}

func (e *codeError) SetId(id string) CodeError {
	e.Id = id
	return e
}

func (e *codeError) SetFailureCode(failureCode int) CodeError {
	e.FailureCode = failureCode
	return e
}

func (e *codeError) GetMeta() MultiMap {
	return e.Meta
}

func (e *codeError) GetStacktrace() (fn string, file string, line int) {
	fn = e.Stacktrace.Fn
	file = e.Stacktrace.File
	line = e.Stacktrace.Line
	return
}

func (e *codeError) Error() string {
	return e.String()
}

func (e *codeError) String() string {
	bb := bytebufferpool.Get()
	defer bytebufferpool.Put(bb)
	_, _ = bb.WriteString("\n")
	if e.Id != "" {
		_, _ = bb.WriteString(fmt.Sprintf("ID      = [%s]\n", e.Id))
	}
	_, _ = bb.WriteString(fmt.Sprintf("CODE    = [%d][%s]\n", e.FailureCode, e.Code))
	_, _ = bb.WriteString(fmt.Sprintf("MESSAGE = %s\n", e.Message))
	if !e.Meta.Empty() {
		_, _ = bb.WriteString("META    = ")
		for i, key := range e.Meta.Keys() {
			values, _ := e.Meta.Values(key)
			if i == 0 {
				_, _ = bb.WriteString(fmt.Sprintf("%s : %v\n", key, values))
			} else {
				_, _ = bb.WriteString(fmt.Sprintf("          %s : %v\n", key, values))
			}
		}
	}
	_, _ = bb.WriteString(fmt.Sprintf("STACK   = %s %s:%d\n", e.Stacktrace.Fn, e.Stacktrace.File, e.Stacktrace.Line))

	return string(bb.Bytes()[:bb.Len()-1])
}

func (e *codeError) ToJson() []byte {
	return jsonEncode(e)
}

func InvalidArgumentError(message string) CodeError {
	return newCodeErrorWithDepth(invalidArgumentErrorFailureCodeCode, invalidArgumentErrorCode, message, 3)
}

func InvalidArgumentErrorWithDetails(message string, details ...string) CodeError {
	err := newCodeErrorWithDepth(invalidArgumentErrorFailureCodeCode, invalidArgumentErrorCode, message, 3)
	if details != nil && len(details) != 0 && len(details)%2 == 0 {
		for i := 0; i < len(details); i = i + 2 {
			k := details[i]
			v := details[i+1]
			err.GetMeta().Add(k, v)
		}
	}
	return err
}

func UnauthorizedError(message string) CodeError {
	return newCodeErrorWithDepth(unauthorizedErrorFailureCodeCode, unauthorizedErrorCode, message, 3)
}

func ForbiddenError(message string) CodeError {
	return newCodeErrorWithDepth(forbiddenErrorFailureCodeCode, forbiddenErrorCode, message, 3)
}

func ForbiddenErrorWithReason(message string, role string, resource ...string) CodeError {
	err := newCodeErrorWithDepth(forbiddenErrorFailureCodeCode, forbiddenErrorCode, message, 3)
	err.GetMeta().Put(role, resource)
	return err
}

func NotFoundError(message string) CodeError {
	return newCodeErrorWithDepth(notFoundErrorFailureCodeCode, notFoundErrorCode, message, 3)
}

func ServiceError(message string) CodeError {
	return newCodeErrorWithDepth(serviceErrorFailureCodeCode, serviceErrorCode, message, 3)
}

func NotImplementedError(message string) CodeError {
	return newCodeErrorWithDepth(serviceNotImplementedErrorFailureCodeCode, serviceNotImplementedErrorCode, message, 3)
}

func UnavailableError(message string) CodeError {
	return newCodeErrorWithDepth(unavailableErrorFailureCodeCode, unavailableErrorCode, message, 3)
}

func NewCodeError(failureCode int, code string, message string) CodeError {
	return newCodeErrorWithDepth(failureCode, code, message, 3)
}

func newCodeErrorWithDepth(failureCode int, code string, message string, skip int) *codeError {
	stacktrace := newStacktrace(skip)
	return &codeError{
		FailureCode: failureCode,
		Code:        code,
		Message:     message,
		Meta:        MultiMap{},
		Stacktrace:  stacktrace,
	}
}

func Transfer(err error) (codeErr CodeError, ok bool) {
	codeErr, ok = err.(CodeError)
	return
}

func FromJson(v []byte) (codeErr CodeError, ok bool) {
	codeErr = &codeError{}
	err := jsonAPI().Unmarshal(v, codeErr)
	if err != nil {
		codeErr = nil
		ok = false
		return
	}
	ok = true
	return
}

func newStacktrace(skip int) stacktrace {
	pc, file, line, ok := runtime.Caller(skip)
	if !ok {
		return stacktrace{
			Fn:   "unknown",
			File: "unknown",
			Line: 0,
		}
	}
	fn := runtime.FuncForPC(pc)
	return stacktrace{
		Fn:   fn.Name(),
		File: fileNameSubGoPath(file),
		Line: line,
	}
}
