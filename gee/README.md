## gin

### 前置知识

go 语言内置了 net/http 库。

其中`http.ListenAndServe(address string, h Handler) error`方法可以监听一个端口，并将所有的 http 请求交给`Handler`的接口实例处理。

查看源码可以发现`Handler`接口实现了一个`ServeHTTP`方法，所以我们只要传入一个实现了该方法的实例即可接手 http 请求的处理。

我们可以通过传入自己编写的`Handler`接口来达到处理 http 请求的目的。

### day1

1. 实现`engine`对象，添加`ServeHTTP`方法和`Router`；
2. `engine`对象实现`Run`方法，内部使用`http.ListenAndServe(address string, h Handler) error`；
