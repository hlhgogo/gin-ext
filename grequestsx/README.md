# grequestsx http请求库

二次封装 grequests组件，可以传递链路信息

## 如何使用

1. 从请求中获取包含链路信息的 context
2. 创建 grequests.RequestOptions 对象，并设置 context
3. 通过 grequestsx 提供的方法发送请求

~~~golang

func GinHandler(c *gin.Context) {
	ctx := c.Request.Context()
	span := tracing.SpanFromContext(ctx) // 从http中提取链路信息
	logger := logrusx.WithContext(ctx)   // 日志信息也将包含链路信息
	if span != nil {
		logger.Info("x-request-id:", span.RequestID())
	}
	grequestsx.Post("https://www.baidu.com/", &grequests.RequestOptions{
		Context: ctx,
	})
	c.JSON(http.StatusOK, gin.H{"result": "success"})
}
~~~

如果无法从链路上获取，可以通过 `tracing.NewContext(context.TODO(),"account-id")` 新建 context

### 启用log

默认不打印请求数据和响应数据，如果需要打印，可以通过定义Flags参数实现：

~~~golang
	got, err := DoRegularRequest(http.MethodGet, "https://postman-echo.com/headers", &grequests.RequestOptions{
		Context: tracing.NewContext(context.TODO(), "xxxxx"),
	}, Flags{EnableLog: true})
~~~