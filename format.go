package errors

import (
	"fmt"
	"github.com/valyala/bytebufferpool"
)

func format(buf *bytebufferpool.ByteBuffer, err CodeError) {
	e := err.(codeError)
	if e.Id() != "" {
		_, _ = buf.WriteString(fmt.Sprintf("ID      = [%s]\n", e.Id()))
	}
	_, _ = buf.WriteString(fmt.Sprintf("CN      = [%d][%s]\n", e.Code(), e.Name()))
	_, _ = buf.WriteString(fmt.Sprintf("MESSAGE = %s\n", e.Message()))
	if len(e.Meta_) > 0 {
		_, _ = buf.WriteString("META    = ")
		metaIdx := 0
		for _, pair := range e.Meta_ {
			if metaIdx == 0 {
				_, _ = buf.WriteString(fmt.Sprintf("%s : %v\n", pair.Key, pair.Value))
			} else {
				_, _ = buf.WriteString(fmt.Sprintf("          %s : %v\n", pair.Key, pair.Value))
			}
			metaIdx++
		}
	}
	fn, file, line := e.Stacktrace()
	_, _ = buf.WriteString(fmt.Sprintf("STACK   = %s %s:%d\n", fn, file, line))
	if e.Cause_ != nil {
		_, _ = buf.WriteString("---\n")
		format(buf, *e.Cause_)
	}
}
