package errors_test

import (
	"fmt"
	"github.com/aacfactory/errors"
	"testing"
)

func TestNewCodeError(t *testing.T) {
	err := errors.NewCodeError(500, "***FOO***", "bar")
	fmt.Println(err)
	fmt.Println(errors.ServiceError("foo"))
	fmt.Println(errors.InvalidArgumentError("foo"))
	fmt.Println(errors.InvalidArgumentErrorWithDetails("foo"))
	fmt.Println(errors.UnauthorizedError("foo"))
	fmt.Println(errors.ForbiddenError("foo"))
	fmt.Println(errors.ForbiddenErrorWithReason("foo", "role", "bar"))
	fmt.Println(errors.NotFoundError("foo"))
	fmt.Println(errors.NotImplementedError("foo"))
	fmt.Println(errors.UnavailableError("foo"))
}

func TestCodeError_ToJson(t *testing.T) {
	fmt.Println(string(errors.ServiceError("x").ToJson()))
}

func TestFromJson(t *testing.T) {
	err := errors.ServiceError("x")
	v := err.ToJson()
	fmt.Println(errors.FromJson(v))
}

func TestTransfer(t *testing.T) {
	var err error = errors.ServiceError("x")
	codeErr, ok := errors.Transfer(err)
	fmt.Println(ok, codeErr.String())
}
