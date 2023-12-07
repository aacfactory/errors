package errors_test

import (
	"bytes"
	"encoding/json"
	serrors "errors"
	"fmt"
	"github.com/aacfactory/errors"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	err := errors.ServiceError("foo").WithCause(errors.ServiceError("bar").WithCause(errors.ServiceError("baz")))
	fmt.Printf("%v\n", err)
	fmt.Printf("%+v\n", err)
}

func TestWrap(t *testing.T) {
	cause := fmt.Errorf("foo")
	err := errors.Wrap(cause)
	fmt.Printf("%v\n", err)
	fmt.Printf("%+v\n", err)
}

func TestCodeError_WithCause(t *testing.T) {
	err := errors.New(500, "***SERVICE EXECUTE FAILED***", "foo")
	err = err.WithCause(fmt.Errorf("bar")).WithCause(fmt.Errorf("baz"))
	fmt.Printf("%v\n", err)
	fmt.Printf("%+v\n", err)
	fmt.Println("contains foo", err.Contains(fmt.Errorf("foo")))
	fmt.Println("contains bar", err.Contains(fmt.Errorf("bar")))
	fmt.Println("contains baz", err.Contains(fmt.Errorf("baz")))
	fmt.Println("contains x  ", err.Contains(fmt.Errorf("x")))
}

func TestCodeError_WithMeta(t *testing.T) {
	err := errors.New(500, "***SERVICE EXECUTE FAILED***", "foo")
	err = err.WithMeta("a", time.Now().String()).WithMeta("b", "b")
	fmt.Printf("%v\n", err)
	fmt.Printf("%+v\n", err)
}

func Test_Json(t *testing.T) {
	err := errors.New(500, "***SERVICE EXECUTE FAILED***", "foo")
	err = err.WithMeta("a", time.Now().String()).WithMeta("b", "b")
	err = err.WithCause(fmt.Errorf("bar")).WithCause(fmt.Errorf("baz"))
	data, _ := json.Marshal(err)
	fmt.Println(string(data))
	err1 := errors.Decode(data)
	fmt.Println(fmt.Sprintf("%+v", err1))
	data1, _ := json.Marshal(err1)
	fmt.Println(bytes.Equal(data, data1))
}

func TestMakeErrors(t *testing.T) {
	errs := errors.MakeErrors()
	for i := 0; i < 3; i++ {
		errs.Append(errors.ServiceError(fmt.Sprintf("%d", i)))
	}
	fmt.Println(fmt.Sprintf("%+v", errs.Error()))
}

func TestContains(t *testing.T) {
	c1 := errors.Warning("c1")
	c2 := errors.Warning("c2").WithCause(c1)
	fmt.Println(errors.Contains(c1, c2), errors.Contains(c2, c1))
	c3 := fmt.Errorf("c3")
	c4 := fmt.Errorf("c4")
	c5 := serrors.Join(c3, c4)
	c6 := errors.Warning("c6").WithCause(c5)
	fmt.Println(errors.Contains(c3, c4))
	fmt.Println(errors.Contains(c5, c3), errors.Contains(c5, c4))
	fmt.Println(errors.Contains(c6, c3))
	fmt.Println(c6.String())
}
