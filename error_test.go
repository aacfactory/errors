package errors_test

import (
	"encoding/json"
	"fmt"
	"github.com/aacfactory/errors"
	"testing"
)

func TestNewCodeError(t *testing.T) {
	err := errors.NewCodeError(500, "***FOO***", "bar")
	err.Stacktrace = append(err.Stacktrace, errors.ErrorStack{
		Fn:   "x",
		File: "x",
		Line: 1,
	})
	err.Meta.Add("a", "a")
	err.Meta.Put("b", nil)
	fmt.Println(err)
	fmt.Println(fmt.Sprintf("xxx %v", err))
	v, _ := json.Marshal(err)
	fmt.Println(string(v))
	//parser.ParseDir()
}
