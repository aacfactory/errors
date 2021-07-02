# 概要
Code Error 是一个继承标准 error 并增加 code 内容进行扩展的错误结构。
# 获取
```go
go get github.com/aacfactory/errors
```
# 使用
```go
// 基本使用
err := errors.NewCodeError(500, "***FOO***", "bar")
fmt.Println(err)
// json 输出
v, _ := json.Marshal(err)
fmt.Println(string(v))

// http status 使用环境
err := errors.InvalidArgumentError("参数错误")
err := errors.InvalidArgumentErrorWithDetails("参数错误", "email", "非法 Email 格式")
err := errors.UnauthorizedError("未认证")
err := errors.ForbiddenError("拒绝访问")
err := errors.ForbiddenErrorWithReason("拒绝访问", "普通用户", "机密文件", "保险箱")
err := errors.NotFoundError("404")
err := errors.ServiceError("服务处理失败")
err := errors.NotImplementedError("功能未实现")
err := errors.UnavailableError("服务不可用")
```