package errors

import (
	"fmt"
	"github.com/rs/xid"
	"github.com/tidwall/gjson"
	"github.com/valyala/bytebufferpool"
	"runtime"
	"strings"
)

const (
	badRequestErrorCode            = 400
	badRequestErrorName            = "***BAD REQUEST***"
	unauthorizedErrorCode          = 401
	unauthorizedErrorName          = "***UNAUTHORIZED***"
	forbiddenErrorCode             = 403
	forbiddenErrorName             = "***FORBIDDEN***"
	notFoundErrorCode              = 404
	notFoundErrorName              = "***NOT FOUND***"
	notAcceptableErrorCode         = 406
	notAcceptableErrorName         = "***NOT ACCEPTABLE***"
	timeoutErrorCode               = 408
	timeoutErrorName               = "***TIMEOUT***"
	serviceErrorCode               = 500
	serviceErrorName               = "***SERVICE EXECUTE FAILED***"
	nilErrorMessage                = "NIL"
	serviceNotImplementedErrorCode = 501
	serviceNotImplementedErrorName = "***SERVICE NOT IMPLEMENTED***"
	unavailableErrorCode           = 503
	unavailableErrorName           = "***SERVICE UNAVAILABLE***"
	warnErrorCode                  = 555
	warnErrorName                  = "***WARNING***"
)

type CodeError interface {
	Id() string
	Code() int
	Name() string
	Message() string
	Stacktrace() (fn string, file string, line int)
	WithMeta(key string, value string) (err CodeError)
	WithCause(cause error) (err CodeError)
	Contains(err error) (has bool)
	Error() string
	Format(state fmt.State, r rune)
	String() string
}

func Empty() CodeError {
	return &codeError{}
}

func BadRequest(message string) CodeError {
	return NewWithDepth(badRequestErrorCode, badRequestErrorName, message, 3)
}

func Unauthorized(message string) CodeError {
	return NewWithDepth(unauthorizedErrorCode, unauthorizedErrorName, message, 3)
}

func Forbidden(message string) CodeError {
	return NewWithDepth(forbiddenErrorCode, forbiddenErrorName, message, 3)
}

func NotFound(message string) CodeError {
	return NewWithDepth(notFoundErrorCode, notFoundErrorName, message, 3)
}

func NotAcceptable(message string) CodeError {
	return NewWithDepth(notAcceptableErrorCode, notAcceptableErrorName, message, 3)
}

func Timeout(message string) CodeError {
	return NewWithDepth(timeoutErrorCode, timeoutErrorName, message, 3)
}

func ServiceError(message string) CodeError {
	return NewWithDepth(serviceErrorCode, serviceErrorName, message, 3)
}

func NilError() CodeError {
	return NewWithDepth(badRequestErrorCode, badRequestErrorName, nilErrorMessage, 3)
}

func NotImplemented(message string) CodeError {
	return NewWithDepth(serviceNotImplementedErrorCode, serviceNotImplementedErrorName, message, 3)
}

func Unavailable(message string) CodeError {
	return NewWithDepth(unavailableErrorCode, unavailableErrorName, message, 3)
}

func Warning(message string) CodeError {
	return NewWithDepth(warnErrorCode, warnErrorName, message, 3)
}

func New(code int, name string, message string) CodeError {
	return NewWithDepth(code, name, message, 3)
}

func NewWithDepth(code int, name string, message string, skip int) CodeError {
	stacktrace_ := newStacktrace(skip)
	return &codeError{
		Id_:         xid.New().String(),
		Code_:       code,
		Name_:       name,
		Message_:    message,
		Meta_:       make(map[string]string),
		Stacktrace_: stacktrace_,
		Cause_:      nil,
	}
}

func Map(err error) (codeErr CodeError) {
	if err == nil {
		codeErr = NewWithDepth(serviceErrorCode, serviceErrorName, "can not map nil to CodeError", 3)
		return
	}
	e, ok := err.(CodeError)
	if ok {
		codeErr = e
		return
	}
	codeErr = NewWithDepth(serviceErrorCode, serviceErrorName, err.Error(), 3)
	return
}

type stacktrace struct {
	Fn   string `json:"fn"`
	File string `json:"file"`
	Line int    `json:"line"`
}

type codeError struct {
	Id_         string            `json:"id,omitempty"`
	Code_       int               `json:"code,omitempty"`
	Name_       string            `json:"name,omitempty"`
	Message_    string            `json:"message,omitempty"`
	Meta_       map[string]string `json:"meta,omitempty"`
	Stacktrace_ stacktrace        `json:"stacktrace,omitempty"`
	Cause_      CodeError         `json:"cause,omitempty"`
}

func (e *codeError) Id() string {
	return e.Id_
}

func (e *codeError) Code() int {
	return e.Code_
}

func (e *codeError) Name() string {
	return e.Name_
}

func (e *codeError) Message() string {
	return e.Message_
}

func (e *codeError) Stacktrace() (fn string, file string, line int) {
	fn = e.Stacktrace_.Fn
	file = e.Stacktrace_.File
	line = e.Stacktrace_.Line
	return
}

func (e *codeError) WithMeta(key string, value string) (err CodeError) {
	e.Meta_[key] = value
	err = e
	return
}

func (e *codeError) WithCause(cause error) (err CodeError) {
	if cause == nil {
		err = e
		return
	}
	ce, ok := cause.(CodeError)
	if !ok {
		ce = NewWithDepth(serviceErrorCode, serviceErrorName, cause.Error(), 4)
	}
	if e.Cause_ == nil {
		e.Cause_ = ce
	} else {
		_ = e.Cause_.WithCause(ce)
	}
	err = e
	return
}

func (e *codeError) Contains(err error) (has bool) {
	if err == nil {
		return
	}
	if e.Message() == err.Error() {
		has = true
		return
	}
	if e.Cause_ != nil {
		has = e.Cause_.Contains(err)
		return
	}
	return
}

func (e *codeError) Error() string {
	return e.String()
}

func (e *codeError) String() string {
	return fmt.Sprintf("%v", e)
}

func (e *codeError) UnmarshalJSON(p []byte) (err error) {
	if p == nil || len(p) == 0 {
		return
	}
	r := gjson.ParseBytes(p)
	if !r.Exists() {
		return
	}

	e.Id_ = r.Get("id").String()
	e.Code_ = int(r.Get("code").Int())
	e.Name_ = r.Get("name").String()
	e.Message_ = r.Get("message").String()
	meta0 := r.Get("meta")
	if meta0.Exists() {
		if e.Meta_ == nil {
			e.Meta_ = make(map[string]string)
		}
		metaValue0 := meta0.Map()
		for key, result := range metaValue0 {
			if result.Exists() {
				e.Meta_[key] = result.String()
			}
		}
	}
	st0 := r.Get("stacktrace")
	if st0.Exists() && st0.IsObject() {
		e.Stacktrace_.File = st0.Get("file").String()
		e.Stacktrace_.Line = int(st0.Get("line").Int())
		e.Stacktrace_.Fn = st0.Get("fn").String()
	}

	cause0 := r.Get("cause")
	if cause0.Exists() && cause0.IsObject() {
		cause := &codeError{}
		causeErr := cause.UnmarshalJSON([]byte(cause0.Raw))
		if causeErr == nil {
			e.Cause_ = cause
		}
	}

	return
}

func (e *codeError) Format(state fmt.State, verb rune) {
	switch verb {
	case 'v':
		switch {
		case state.Flag('+'):
			buf := bytebufferpool.Get()
			_, _ = buf.WriteString(">>>>>>>>>>>>>\n")
			if e.Id() != "" {
				_, _ = buf.WriteString(fmt.Sprintf("ID      = [%s]\n", e.Id()))
			}
			_, _ = buf.WriteString(fmt.Sprintf("CN      = [%d][%s]\n", e.Code(), e.Name()))
			_, _ = buf.WriteString(fmt.Sprintf("MESSAGE = %s\n", e.Message()))
			if len(e.Meta_) > 0 {
				_, _ = buf.WriteString("META    = ")
				metaIdx := 0
				for k, v := range e.Meta_ {
					if metaIdx == 0 {
						_, _ = buf.WriteString(fmt.Sprintf("%s : %v\n", k, v))
					} else {
						_, _ = buf.WriteString(fmt.Sprintf("          %s : %v\n", k, v))
					}
					metaIdx++
				}
			}
			fn, file, line := e.Stacktrace()
			_, _ = buf.WriteString(fmt.Sprintf("STACK   = %s %s:%d\n", fn, file, line))
			formatCause(buf, e.Cause_, 0)
			_, _ = buf.WriteString("<<<<<<<<<<<<<\n")
			content := buf.Bytes()[:buf.Len()-1]
			bytebufferpool.Put(buf)
			_, _ = fmt.Fprintf(state, "%s", content)
		default:
			_, _ = fmt.Fprintf(state, "%s", e.Message())
		}
	default:
		_, _ = fmt.Fprintf(state, "%s", e.Message())
	}
}

func formatCause(buf *bytebufferpool.ByteBuffer, cause CodeError, depth int) {
	if cause == nil {
		return
	}
	if depth == 0 {
		_, _ = buf.WriteString(fmt.Sprintf("CAUSE   = %s\n", cause.Message()))
	} else {
		_, _ = buf.WriteString(fmt.Sprintf("        = %s\n", cause.Message()))
	}
	e, ok := cause.(*codeError)
	if !ok {
		return
	}
	if e.Cause_ != nil {
		formatCause(buf, e.Cause_, depth+1)
	}
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
	if strings.IndexByte(file, '/') == 0 || strings.IndexByte(file, ':') == 1 {
		idx := strings.Index(file, "/src/")
		if idx > 0 {
			file = file[idx+5:]
		} else {
			idx = strings.Index(file, "/pkg/mod/")
			if idx > 0 {
				file = file[idx+9:]
			}
		}
	}
	fn := runtime.FuncForPC(pc)
	return stacktrace{
		Fn:   fn.Name(),
		File: file,
		Line: line,
	}
}
