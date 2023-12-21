package errors

import (
	"github.com/valyala/bytebufferpool"
	"sort"
	"unsafe"
)

type Pair struct {
	Key   string `json:"key" avro:"key"`
	Value string `json:"value" avro:"value"`
}

func (pair Pair) MarshalJSON() (p []byte, err error) {
	buf := bytebufferpool.Get()
	_, _ = buf.Write(lb)
	_, _ = buf.Write(dqm)
	_, _ = buf.Write(keyIdent)
	_, _ = buf.Write(dqm)
	_, _ = buf.Write(colon)
	_, _ = buf.Write(dqm)
	_, _ = buf.WriteString(pair.Key)
	_, _ = buf.Write(dqm)
	_, _ = buf.Write(comma)
	_, _ = buf.Write(dqm)
	_, _ = buf.Write(valueIdent)
	_, _ = buf.Write(dqm)
	_, _ = buf.Write(colon)
	_, _ = buf.Write(dqm)
	_, _ = buf.WriteString(pair.Value)
	_, _ = buf.Write(dqm)
	_, _ = buf.Write(rb)
	s := buf.String()
	p = unsafe.Slice(unsafe.StringData(s), len(s))
	bytebufferpool.Put(buf)
	return
}

type Meta []Pair

func (m Meta) Len() int {
	return len(m)
}

func (m Meta) Less(i, j int) bool {
	return m[i].Key < m[j].Key
}

func (m Meta) Swap(i, j int) {
	m[i], m[j] = m[j], m[i]
}

func (m Meta) Add(key string, value string) Meta {
	for i, pair := range m {
		if pair.Key == key {
			pair.Value = value
			m[i] = pair
			return m
		}
	}
	n := append(m, Pair{
		Key:   key,
		Value: value,
	})
	sort.Sort(n)
	return n
}

func (m Meta) MarshalJSON() (p []byte, err error) {
	buf := bytebufferpool.Get()
	_, _ = buf.Write(lqb)
	if m.Len() == 0 {
		_, _ = buf.Write(rqb)
		p = buf.Bytes()
		bytebufferpool.Put(buf)
		return
	}
	for i, pair := range m {
		if i > 0 {
			_, _ = buf.Write(comma)
		}
		b, _ := pair.MarshalJSON()
		_, _ = buf.Write(b)
	}
	_, _ = buf.Write(rqb)
	s := buf.String()
	p = unsafe.Slice(unsafe.StringData(s), len(s))
	bytebufferpool.Put(buf)
	return
}
