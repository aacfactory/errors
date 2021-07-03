package errors

import (
	"sync"
)

var _once = new(sync.Once)

func init()  {
	_once.Do(func() {
		initEnv()
		initJsonApi()
	})
}