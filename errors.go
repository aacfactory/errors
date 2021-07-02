package errors

import (
	"fmt"
	"github.com/pharosnet/errors"
	"github.com/tidwall/gjson"
	"github.com/valyala/bytebufferpool"
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

type ErrorStack struct {
	Fn   string `json:"fn"`
	File string `json:"file"`
	Line int    `json:"line"`
}

type CodeError struct {
	Id          string       `json:"id,omitempty"`
	FailureCode int          `json:"failureCode,omitempty"`
	Code        string       `json:"code,omitempty"`
	Message     string       `json:"message,omitempty"`
	Meta        MultiMap     `json:"meta,omitempty"`
	Stacktrace  []ErrorStack `json:"stacktrace,omitempty"`
}

func (e *CodeError) SetId(id string) *CodeError {
	e.Id = id
	return e
}

func (e *CodeError) SetFailureCode(failureCode int) *CodeError {
	e.FailureCode = failureCode
	return e
}

func (e *CodeError) Error() string {
	return e.String()
}

func (e *CodeError) String() string {
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
	_, _ = bb.WriteString("STACK   = ")
	for i, stack := range e.Stacktrace {
		if i == 0 {
			_, _ = bb.WriteString(fmt.Sprintf("%s %s:%d\n", stack.Fn, stack.File, stack.Line))
		} else {
			_, _ = bb.WriteString(fmt.Sprintf("          %s %s:%d\n", stack.Fn, stack.File, stack.Line))
		}
	}
	return string(bb.Bytes()[:bb.Len()-1])
}

func InvalidArgumentError(message string) *CodeError {
	return NewCodeError(invalidArgumentErrorFailureCodeCode, invalidArgumentErrorCode, message)
}

func InvalidArgumentErrorWithDetails(message string, details ...string) *CodeError {
	err := NewCodeError(invalidArgumentErrorFailureCodeCode, invalidArgumentErrorCode, message)
	if details != nil && len(details) != 0 && len(details)%2 == 0 {
		for i := 0; i < len(details); i = i + 2 {
			k := details[i]
			v := details[i+1]
			err.Meta.Add(k, v)
		}
	}
	return err
}

func UnauthorizedError(message string) *CodeError {
	return NewCodeError(unauthorizedErrorFailureCodeCode, unauthorizedErrorCode, message)
}

func ForbiddenError(message string) *CodeError {
	return NewCodeError(forbiddenErrorFailureCodeCode, forbiddenErrorCode, message)
}

func ForbiddenErrorWithReason(message string, role string, resource ...string) *CodeError {
	err := NewCodeError(forbiddenErrorFailureCodeCode, forbiddenErrorCode, message)
	err.Meta.Put(role, resource)
	return err
}

func NotFoundError(message string) *CodeError {
	return NewCodeError(notFoundErrorFailureCodeCode, notFoundErrorCode, message)
}

func ServiceError(message string) *CodeError {
	return NewCodeError(serviceErrorFailureCodeCode, serviceErrorCode, message)
}

func NotImplementedError(message string) *CodeError {
	return NewCodeError(serviceNotImplementedErrorFailureCodeCode, serviceNotImplementedErrorCode, message)
}

func UnavailableError(message string) *CodeError {
	return NewCodeError(unavailableErrorFailureCodeCode, unavailableErrorCode, message)
}

func NewCodeError(failureCode int, code string, message string) *CodeError {
	err := errors.NewWithDepth(1, 4, message)
	stacktrace := make([]ErrorStack, 0, 1)
	stackJsonField := gjson.Get(fmt.Sprintf("%-v", err), "stack")
	if stackJsonField.Exists() {
		jsonDecodeFromString(stackJsonField.String(), &stacktrace)
	}
	return &CodeError{
		FailureCode: failureCode,
		Code:        code,
		Message:     message,
		Meta:        MultiMap{},
		Stacktrace:  stacktrace,
	}
}

func NewCodeErrorWithCause(failureCode int, code string, message string, cause error) *CodeError {
	err := errors.WithDepth(1, 4, cause, message)
	stacktrace := make([]ErrorStack, 0, 1)
	stackJsonField := gjson.Get(fmt.Sprintf("%-v", err), "stack")
	if stackJsonField.Exists() {
		jsonDecodeFromString(stackJsonField.String(), &stacktrace)
	}
	return &CodeError{
		FailureCode: failureCode,
		Code:        code,
		Message:     message,
		Meta:        MultiMap{},
		Stacktrace:  stacktrace,
	}
}

func NewCodeErrorWithDepth(failureCode int, code string, message string, depth int) *CodeError {
	err := errors.NewWithDepth(depth, 3+depth, message)
	stacktrace := make([]ErrorStack, 0, 1)
	stackJsonField := gjson.Get(fmt.Sprintf("%-v", err), "stack")
	if stackJsonField.Exists() {
		jsonDecodeFromString(stackJsonField.String(), &stacktrace)
	}
	return &CodeError{
		FailureCode: failureCode,
		Code:        code,
		Message:     message,
		Meta:        MultiMap{},
		Stacktrace:  stacktrace,
	}
}

