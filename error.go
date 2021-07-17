package errors

import (
	"fmt"
	"github.com/valyala/bytebufferpool"
	"runtime"
)

const (
	invalidArgumentErrorFailureCodeCode       = 400
	invalidArgumentErrorCode                  = "***BAD REQUEST***"
	unauthorizedErrorFailureCodeCode          = 401
	unauthorizedErrorCode                     = "***UNAUTHORIZED***"
	forbiddenErrorFailureCodeCode             = 403
	forbiddenErrorCode                        = "***FORBIDDEN***"
	notFoundErrorFailureCodeCode              = 404
	notFoundErrorCode                         = "***NOT FOUND***"
	serviceErrorFailureCodeCode               = 500
	serviceErrorCode                          = "***SERVICE EXECUTE FAILED***"
	serviceNotImplementedErrorFailureCodeCode = 501
	serviceNotImplementedErrorCode            = "***SERVICE NOT IMPLEMENTED***"
	unavailableErrorFailureCodeCode           = 503
	unavailableErrorCode                      = "***SERVICE UNAVAILABLE***"
)

type CodeError interface {
	Id() string
	Code() string
	FailureCode() int
	Message() string
	Meta() MultiMap
	Stacktrace() (fn string, file string, line int)
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
	Id_          string     `json:"id,omitempty"`
	FailureCode_ int        `json:"failureCode,omitempty"`
	Code_        string     `json:"code,omitempty"`
	Message_     string     `json:"message,omitempty"`
	Meta_        MultiMap   `json:"meta,omitempty"`
	Stacktrace_  stacktrace `json:"stacktrace,omitempty"`
}

func (e *codeError) Id() string {
	return e.Id_
}

func (e *codeError) Code() string {
	return e.Code_
}

func (e *codeError) FailureCode() int {
	return e.FailureCode_
}

func (e *codeError) Message() string {
	return e.Message_
}

func (e *codeError) Meta() MultiMap {
	return e.Meta_
}

func (e *codeError) Stacktrace() (fn string, file string, line int) {
	fn = e.Stacktrace_.Fn
	file = e.Stacktrace_.File
	line = e.Stacktrace_.Line
	return
}

func (e *codeError) Error() string {
	return e.String()
}

func (e *codeError) String() string {
	bb := bytebufferpool.Get()
	defer bytebufferpool.Put(bb)
	_, _ = bb.WriteString("\n")
	if e.Id() != "" {
		_, _ = bb.WriteString(fmt.Sprintf("ID      = [%s]\n", e.Id()))
	}
	_, _ = bb.WriteString(fmt.Sprintf("CODE    = [%d][%s]\n", e.FailureCode(), e.Code()))
	_, _ = bb.WriteString(fmt.Sprintf("MESSAGE = %s\n", e.Message()))
	if !e.Meta().Empty() {
		_, _ = bb.WriteString("META    = ")
		for i, key := range e.Meta().Keys() {
			values, _ := e.Meta().Values(key)
			if i == 0 {
				_, _ = bb.WriteString(fmt.Sprintf("%s : %v\n", key, values))
			} else {
				_, _ = bb.WriteString(fmt.Sprintf("          %s : %v\n", key, values))
			}
		}
	}
	fn, file, line := e.Stacktrace()
	_, _ = bb.WriteString(fmt.Sprintf("STACK   = %s %s:%d\n", fn, file, line))

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
			err.Meta().Add(k, v)
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
	err.Meta().Put(role, resource)
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
		FailureCode_: failureCode,
		Code_:        code,
		Message_:     message,
		Meta_:        MultiMap{},
		Stacktrace_:  stacktrace,
	}
}

func Transfer(err error) (codeErr CodeError, ok bool) {
	codeErr, ok = err.(CodeError)
	return
}

func Map(err error) (codeErr CodeError) {
	ok := false
	codeErr, ok = Transfer(err)
	if ok {
		return
	}
	codeErr = ServiceError(err.Error())
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
