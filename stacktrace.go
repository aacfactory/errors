package errors

import (
	"github.com/valyala/bytebufferpool"
	"runtime"
	"strconv"
	"strings"
	"unsafe"
)

func newStacktrace(skip int) Stacktrace {
	pc, file, line, ok := runtime.Caller(skip)
	if !ok {
		return Stacktrace{
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
	return Stacktrace{
		Fn:   fn.Name(),
		File: file,
		Line: line,
	}
}

type Stacktrace struct {
	Fn   string `json:"fn" avro:"fn"`
	File string `json:"file" avro:"file"`
	Line int    `json:"line" avro:"line"`
}

func (s Stacktrace) MarshalJSON() (p []byte, err error) {
	buf := bytebufferpool.Get()
	_, _ = buf.Write(lb)
	// fn
	_, _ = buf.Write(dqm)
	_, _ = buf.Write(fnIdent)
	_, _ = buf.Write(dqm)
	_, _ = buf.Write(colon)
	_, _ = buf.Write(dqm)
	_, _ = buf.WriteString(s.Fn)
	_, _ = buf.Write(dqm)
	_, _ = buf.Write(comma)
	// file
	_, _ = buf.Write(dqm)
	_, _ = buf.Write(fileIdent)
	_, _ = buf.Write(dqm)
	_, _ = buf.Write(colon)
	_, _ = buf.Write(dqm)
	_, _ = buf.WriteString(s.File)
	_, _ = buf.Write(dqm)
	_, _ = buf.Write(comma)
	// lind
	_, _ = buf.Write(dqm)
	_, _ = buf.Write(lineIdent)
	_, _ = buf.Write(dqm)
	_, _ = buf.Write(colon)
	_, _ = buf.WriteString(strconv.Itoa(s.Line))
	_, _ = buf.Write(rb)
	ss := buf.String()
	p = unsafe.Slice(unsafe.StringData(ss), len(ss))
	bytebufferpool.Put(buf)
	return
}
