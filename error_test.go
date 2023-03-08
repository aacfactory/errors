package errors_test

import (
	"bytes"
	"encoding/json"
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

func TestMap(t *testing.T) {
	cause := fmt.Errorf("foo")
	err := errors.Map(cause)
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
