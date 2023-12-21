package errors

import (
	"github.com/valyala/bytebufferpool"
	"strconv"
	"unsafe"
)

var (
	lb              = []byte{'{'}
	rb              = []byte{'}'}
	lqb             = []byte{'['}
	rqb             = []byte{']'}
	dqm             = []byte{'"'}
	comma           = []byte{','}
	colon           = []byte{':'}
	idIdent         = []byte("id")
	codeIdent       = []byte("code")
	nameIdent       = []byte("name")
	messageIdent    = []byte("message")
	metaIdent       = []byte("meta")
	keyIdent        = []byte("key")
	valueIdent      = []byte("value")
	stacktraceIdent = []byte("stacktrace")
	fnIdent         = []byte("fn")
	fileIdent       = []byte("file")
	lineIdent       = []byte("line")
	causeIdent      = []byte("cause")
)

func (e CodeErrorImpl) MarshalJSON() (p []byte, err error) {
	buf := bytebufferpool.Get()
	_, _ = buf.Write(lb)
	// id
	_, _ = buf.Write(dqm)
	_, _ = buf.Write(idIdent)
	_, _ = buf.Write(dqm)
	_, _ = buf.Write(colon)
	_, _ = buf.Write(dqm)
	_, _ = buf.WriteString(e.Id_)
	_, _ = buf.Write(dqm)
	_, _ = buf.Write(comma)
	// code
	_, _ = buf.Write(dqm)
	_, _ = buf.Write(codeIdent)
	_, _ = buf.Write(dqm)
	_, _ = buf.Write(colon)
	_, _ = buf.WriteString(strconv.Itoa(e.Code_))
	_, _ = buf.Write(comma)
	// name
	_, _ = buf.Write(dqm)
	_, _ = buf.Write(nameIdent)
	_, _ = buf.Write(dqm)
	_, _ = buf.Write(colon)
	_, _ = buf.Write(dqm)
	_, _ = buf.WriteString(e.Name_)
	_, _ = buf.Write(dqm)
	_, _ = buf.Write(comma)
	// message
	_, _ = buf.Write(dqm)
	_, _ = buf.Write(messageIdent)
	_, _ = buf.Write(dqm)
	_, _ = buf.Write(colon)
	_, _ = buf.Write(dqm)
	_, _ = buf.WriteString(e.Message_)
	_, _ = buf.Write(dqm)
	_, _ = buf.Write(comma)
	// Meta
	_, _ = buf.Write(dqm)
	_, _ = buf.Write(metaIdent)
	_, _ = buf.Write(dqm)
	_, _ = buf.Write(colon)
	metaBytes, _ := e.Meta_.MarshalJSON()
	_, _ = buf.Write(metaBytes)
	_, _ = buf.Write(comma)
	// Stacktrace
	_, _ = buf.Write(dqm)
	_, _ = buf.Write(stacktraceIdent)
	_, _ = buf.Write(dqm)
	_, _ = buf.Write(colon)
	stacktraceBytes, _ := e.Stacktrace_.MarshalJSON()
	_, _ = buf.Write(stacktraceBytes)
	// cause
	if e.Cause_ != nil {
		_, _ = buf.Write(comma)
		_, _ = buf.Write(dqm)
		_, _ = buf.Write(causeIdent)
		_, _ = buf.Write(dqm)
		_, _ = buf.Write(colon)
		causeBytes, _ := e.Cause_.MarshalJSON()
		_, _ = buf.Write(causeBytes)
	}
	_, _ = buf.Write(rb)
	s := buf.String()
	p = unsafe.Slice(unsafe.StringData(s), len(s))
	bytebufferpool.Put(buf)
	return
}
