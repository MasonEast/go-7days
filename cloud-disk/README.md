# CloudDisk

> 轻量级云盘系统，基于 go-zero、xorm 实现。

使用到的命令

```text
# 创建API服务
goctl api new core
# 启动服务
go run core.go -f etc/core-api.yaml
# 使用api文件生成代码
goctl api go -api core.api -dir . -style go_zero
```

## go-zero 开发流程

1. 定义 api
2. 使用 api 命令生成代码

## 关键功能

- 文件上传：

  - 从请求读取文件；
  - 截取文件内容存 hash，方便后续比对是否资源重复；
  - 与数据库中 hash 比对，查重；
  - 上传 cos 获取文件上传路径；
  - 将文件相关信息落库；

- 鉴权：
  - 检测用户请求头是否有 x-token；
  - 对 x-token 内容进行解析，判断是否有效；
  - 将解析内容存到请求头中，方便后续使用；

腾讯云 COS 后台地址：https://console.cloud.tencent.com/cos/bucket

腾讯云 COS 帮助文档：https://cloud.tencent.com/document/product/436/31215
