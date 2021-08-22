# 概要

Code Error 是一个继承标准 error 并增加 code 内容进行扩展的错误结构。

# 获取

```go
go get github.com/aacfactory/errors
```

# 使用

```go
// 基本使用
err := errors.New(500, errors.DefaultErrorName, "foo")
fmt.Printf("%v\n", err)
// Output: foo
fmt.Printf("%+v\n", err)
// Output: 
/*
   >>>>>>>>>>>>>
   ID      = [c4h3pfde2f60qd5fde8g]
   CN      = [500][***SERVICE ERROR***]
   MESSAGE = foo
   STACK   = github.com/aacfactory/errors_test.TestNew github.com/aacfactory/errors/error_test.go:11
   <<<<<<<<<<<<<
*/
```

```go
// Map (标准 error 转换 CodeError)
cause := fmt.Errorf("foo")
err := errors.Map(cause)
fmt.Printf("%v\n", err)
fmt.Printf("%+v\n", err)
```

```go
// WithCause (引入其它错误)
err := errors.New(500, errors.DefaultErrorName, "foo")
err = err.WithCause(fmt.Errorf("bar")).WithCause(fmt.Errorf("baz"))
fmt.Printf("%v\n", err)
fmt.Printf("%+v\n", err)
fmt.Println("contains foo", err.Contains(fmt.Errorf("foo"))) // true
fmt.Println("contains bar", err.Contains(fmt.Errorf("bar"))) // true
fmt.Println("contains baz", err.Contains(fmt.Errorf("baz"))) // true
fmt.Println("contains x  ", err.Contains(fmt.Errorf("x"))) // false
```

```go
// WithMeta (添加员元数据，一般用于参数校验错误)
err := errors.New(500, errors.DefaultErrorName, "foo")
err = err.WithMeta("a", time.Now().String()).WithMeta("b", "b")
fmt.Printf("%v\n", err)
fmt.Printf("%+v\n", err)
```

## Json 输出

```json
{
  "id": "c4h3rsde2f60j14h822g",
  "code": 500,
  "name": "***SERVICE ERROR***",
  "message": "foo",
  "meta": {
    "a": "2021-08-22 20:07:13.0526095 +0800 CST m=+0.002060201",
    "b": "b"
  },
  "stacktrace": {
    "fn": "github.com/aacfactory/errors_test.Test_Json",
    "file": "github.com/aacfactory/errors/error_test.go",
    "line": 43
  },
  "cause": {
    "id": "c4h3rsde2f60j14h8230",
    "code": 500,
    "name": "***SERVICE ERROR***",
    "message": "bar",
    "stacktrace": {
      "fn": "testing.tRunner",
      "file": "testing/testing.go",
      "line": 1193
    },
    "cause": {
      "id": "c4h3rsde2f60j14h823g",
      "code": 500,
      "name": "***SERVICE ERROR***",
      "message": "baz",
      "stacktrace": {
        "fn": "testing.tRunner",
        "file": "testing/testing.go",
        "line": 1193
      }
    }
  }
}
```

## 引用感谢

* [valyala/bytebufferpool](https://github.com/valyala/bytebufferpool)
* [rs/xid](https://github.com/rs/xid)
