## gin

### 前置知识

go 语言内置了 net/http 库。

其中`http.ListenAndServe(address string, h Handler) error`方法可以监听一个端口，并将所有的 http 请求交给`Handler`的接口实例处理。

查看源码可以发现`Handler`接口实现了一个`ServeHTTP`方法，所以我们只要传入一个实现了该方法的实例即可接手 http 请求的处理。

我们可以通过传入自己编写的`Handler`接口来达到处理 http 请求的目的。

### day1

1. 实现`engine`对象，添加`ServeHTTP`方法和`Router`；
2. `engine`对象实现`Run`方法，内部使用`http.ListenAndServe(address string, h Handler) error`；

### day2

#### 设计 Context

1. 对 web 服务来说，就是根据请求\*http.Request，构造响应 http.ResponseWriter。但是这两个对象提供的接口粒度太细，会导致用户写大量重复代码。
2. Context 随着每个请求的出现而产生，请求结束而销毁，和当前请求强相关的信息都应由 Context 来承载。
