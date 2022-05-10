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

腾讯云 COS 后台地址：https://console.cloud.tencent.com/cos/bucket

腾讯云 COS 帮助文档：https://cloud.tencent.com/document/product/436/31215
