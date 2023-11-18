# douyin-security

抖音开放平台内容安全检测包括文本、图片

# Install

```shell
go get github.com/xiaoxuan6/douyin-security
```

# Usage

```go
// 创建账号
account := NewAccount("xxx", "xxx", true) // app_id, app_secret, 是否为沙箱环境
client := NewClient(account)

// 获取access_token，有效期为2小时，自行缓存，如何不设置每次请求都会获取
accessToken, err := client.GetAccessToken()
client.SetAccessToken(accessToken)

// 设置最大请求次数，默认为3次
client.SetMaxAttempts(3)

// 文本检测
response := client.TextVerify("测试文本")

// 图片检测
response := client.ImageVerify("http://xxxxx.jpg")
or
response := client.ImageVerify("file_base64")
```



