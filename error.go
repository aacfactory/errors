package errors

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/rs/xid"
	"github.com/valyala/bytebufferpool"
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
	tooEarlyCode                   = 425
	tooEarlyName                   = "***TOO EARLY***"
	tooManyRequestsCode            = 429
	tooManyRequestsName            = "***TOO MANY REQUEST***"
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
	json.Marshaler
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
	return NewWithDepth(notFoundErrorCode, notFoundErrorName, nilErrorMessage, 3)
}

func NotImplemented(message string) CodeError {
	return NewWithDepth(serviceNotImplementedErrorCode, serviceNotImplementedErrorName, message, 3)
}

func Unavailable(message string) CodeError {
	return NewWithDepth(unavailableErrorCode, unavailableErrorName, message, 3)
}

func TooMayRequest(message string) CodeError {
	return NewWithDepth(tooManyRequestsCode, tooManyRequestsName, message, 3)
}

func TooEarly(message string) CodeError {
	return NewWithDepth(tooEarlyCode, tooEarlyName, message, 3)
}

func Warning(message string) CodeError {
	return NewWithDepth(warnErrorCode, warnErrorName, message, 3)
}

func New(code int, name string, message string) CodeError {
	return NewWithDepth(code, name, message, 3)
}

func NewWithDepth(code int, name string, message string, skip int) CodeError {
	return codeError{
		Id_:         xid.New().String(),
		Code_:       code,
		Name_:       name,
		Message_:    message,
		Meta_:       nil,
		Stacktrace_: newStacktrace(skip),
		Cause_:      nil,
	}
}

func Wrap(err error) (codeErr CodeError) {
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

func Decode(p []byte) (err CodeError) {
	v := codeError{}
	decodeErr := json.Unmarshal(p, &v)
	if decodeErr != nil {
		err = Warning("decode code error failed").WithCause(decodeErr)
		return
	}
	err = v
	return
}

type codeError struct {
	Id_         string     `json:"id,omitempty"`
	Code_       int        `json:"code,omitempty"`
	Name_       string     `json:"name,omitempty"`
	Message_    string     `json:"message,omitempty"`
	Meta_       meta       `json:"meta,omitempty"`
	Stacktrace_ stacktrace `json:"stacktrace,omitempty"`
	Cause_      *codeError `json:"cause,omitempty"`
}

func (e codeError) Id() string {
	return e.Id_
}

func (e codeError) Code() int {
	return e.Code_
}

func (e codeError) Name() string {
	return e.Name_
}

func (e codeError) Message() string {
	return e.Message_
}

func (e codeError) Stacktrace() (fn string, file string, line int) {
	fn = e.Stacktrace_.Fn
	file = e.Stacktrace_.File
	line = e.Stacktrace_.Line
	return
}

func (e codeError) WithMeta(key string, value string) (err CodeError) {
	e.Meta_ = e.Meta_.Add(key, value)
	err = e
	return
}

func (e codeError) WithCause(cause error) (err CodeError) {
	if cause == nil {
		err = e
		return
	}
	ce, ok := cause.(CodeError)
	if !ok {
		joined, isJoined := cause.(JoinedErrors)
		if isJoined {
			errs := joined.Unwrap()
			if len(errs) > 0 {
				ce = NewWithDepth(serviceErrorCode, serviceErrorName, errs[0].Error(), 4)
				for _, sub := range errs[1:] {
					ce = ce.WithCause(sub)
				}
			}
		} else {
			ce = NewWithDepth(serviceErrorCode, serviceErrorName, cause.Error(), 4)
		}
	}
	if e.Cause_ == nil {
		ca := ce.(codeError)
		e.Cause_ = &ca
	} else {
		ce = e.Cause_.WithCause(ce)
		ca := ce.(codeError)
		e.Cause_ = &ca
	}
	err = e
	return
}

func (e codeError) Contains(err error) (has bool) {
	if err == nil {
		return
	}
	codeErr, ok := err.(CodeError)
	if ok {
		if e.Message() == codeErr.Message() {
			has = true
		}
	} else {
		if e.Message() == err.Error() {
			has = true
		} else {
			has = errors.Is(e, err)
		}
	}
	if !has && e.Cause_ != nil {
		has = e.Cause_.Contains(err)
	}
	return
}

func (e codeError) Error() string {
	return e.String()
}

func (e codeError) String() string {
	return fmt.Sprintf("%+v", e)
}

func (e codeError) Format(state fmt.State, verb rune) {
	switch verb {
	case 'v':
		switch {
		case state.Flag('+'):
			buf := bytebufferpool.Get()
			_, _ = buf.WriteString("\n>>>>>>>>>>>>>\n")
			format(buf, e)
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

func MakeErrors() Errors {
	return make([]CodeError, 0, 1)
}

type Errors []CodeError

func (e *Errors) Append(err error) {
	*e = append(*e, Wrap(err))
}

func (e *Errors) Error() (err error) {
	if len(*e) == 0 {
		return
	}
	e0 := (*e)[0]
	if len(*e) > 1 {
		for i := 1; i < len(*e); i++ {
			e0 = e0.WithCause((*e)[i])
		}
	}
	err = e0
	return
}

func Contains(a error, b error) (has bool) {
	if a == nil {
		return
	}
	if b == nil {
		return
	}
	codeErrorA, aOk := a.(CodeError)
	if aOk {
		has = codeErrorA.Contains(b)
		return
	}
	has = errors.Is(a, b)
	return
}

func As(err error) (e CodeError, ok bool) {
	if err == nil {
		return
	}
	e, ok = err.(CodeError)
	return
}

type JoinedErrors interface {
	Unwrap() []error
}
