# 巡汐 Probelune 协议

[English](README.md) · [产品文档与 PRD](https://github.com/sihuangtech/probelune)

跨端复用的 Go 数据结构、参数校验、网络访问策略与 JSON 协议说明。**它不是独立运行的后端，也不需要部署。**

需要 Go 1.27 或更新版本：

```sh
go get github.com/sihuangtech/probelune-contracts@v0.1.1
```

- 面板后端和探测节点直接引用带版本号的 Go 模块。
- 网页、桌面和 Flutter 按相同 JSON 协议通信，各自维护适合本端的模型。
- `models.go` 定义监测、诊断、节点租约、结果与版本协商。
- `validate.go` 校验参数并填充默认值，执行器入口再次校验。
- `network.go` 校验解析 IP、拦截默认禁止的私网地址，并固定实际拨号 IP。
- 各应用独立仓库、独立发布，不依赖相邻源码目录。

协议 v1、模块 v0.1.0。兼容字段可增量添加；破坏性变更必须升级协议主版本。

文档：[节点协议](docs/PROTOCOL.md)、[HTTP API](docs/API.md)、[JSON Schema](schemas/check.schema.json)、[OpenAPI](openapi.yaml)。

开发检查：`go fmt ./...`、`go vet ./...`、`go test ./...`。协议说明以英文为主，代码含必要中文注释。贡献与安全报告方式见 [CONTRIBUTING.md](CONTRIBUTING.md)、[SECURITY.md](SECURITY.md)。许可证为 GPL-3.0-or-later。
