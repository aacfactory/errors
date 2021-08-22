package errors

import (
	"fmt"
	"github.com/rs/xid"
	"github.com/valyala/bytebufferpool"
	"runtime"
	"strings"
)

const (
	DefaultErrorCode = 500
	DefaultErrorName = "***SERVICE ERROR***"
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
		Cause:       nil,
	}
}

func Map(err error) (codeErr CodeError) {
	if err == nil {
		codeErr = NewWithDepth(DefaultErrorCode, DefaultErrorName, "can not map nil to CodeError", 3)
		return
	}
	e, ok := err.(CodeError)
	if ok {
		codeErr = e
		return
	}
	codeErr = NewWithDepth(DefaultErrorCode, DefaultErrorName, err.Error(), 3)
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
	Cause       CodeError         `json:"cause,omitempty"`
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
		ce = NewWithDepth(DefaultErrorCode, DefaultErrorName, cause.Error(), 4)
	}
	if e.Cause == nil {
		e.Cause = ce
	} else {
		_ = e.Cause.WithCause(ce)
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
	if e.Cause != nil {
		has = e.Cause.Contains(err)
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
			formatCause(buf, e.Cause, 0)
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
	if e.Cause != nil {
		formatCause(buf, e.Cause, depth+1)
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
		}
	}
	fn := runtime.FuncForPC(pc)
	return stacktrace{
		Fn:   fn.Name(),
		File: file,
		Line: line,
	}
}
