# 巡汐 Probelune 协议

[English](README.md) · [产品文档与 PRD](https://github.com/sihuangtech/probelune)

跨端复用的 Go 数据结构、参数校验、网络访问策略与 JSON 协议说明。**它不是独立运行的后端，也不需要部署。**

需要 Go 1.27 或更新版本：

```sh
export GONOPROXY=github.com/sihuangtech/probelune-contracts
export GONOSUMDB=github.com/sihuangtech/probelune-contracts
export GOMODCACHE="$PWD/.cache/go-mod-v0.1.1"
go get github.com/sihuangtech/probelune-contracts@v0.1.1
```

所有仓库统一使用 `v0.1.1`，包含证书剩余有效期检查字段。
- 网页、桌面和 Flutter 按相同 JSON 协议通信，各自维护适合本端的模型。
- `check.go` 定义检测类型、配置与能力要求；`models.go` 定义监测、诊断、节点租约、结果与版本协商。
- `validate.go` 校验参数并填充默认值，执行器入口再次校验。
- `network.go` 校验解析 IP、拦截默认禁止的私网地址，并固定实际拨号 IP。
- 各应用独立仓库、独立发布，不依赖相邻源码目录。

协议 v1、模块 v0.1.1。兼容字段可增量添加；破坏性变更必须升级协议主版本。

发布记录统一维护在 [GitHub Releases](https://github.com/sihuangtech/probelune-contracts/releases)。

文档：[节点协议](docs/PROTOCOL.md)、[HTTP API](docs/API.md)、[JSON Schema](schemas/check.schema.json)、[OpenAPI](openapi.yaml)。

开发检查：`go fmt ./...`、`go vet ./...`、`go test ./...`。协议说明以英文为主，代码含必要中文注释。贡献与安全报告方式见 [CONTRIBUTING.md](CONTRIBUTING.md)、[SECURITY.md](SECURITY.md)。许可证为 GPL-3.0-or-later。


## 了解 Probelune

协议库为面板、探测节点和客户端之间的通信提供统一的数据定义与校验规则。

开始自部署时，先配置网页面板和探测节点，再按需要连接手机或桌面客户端。完整架构与部署指南可以在项目总览中找到。

| 组件 | 用途 |
| --- | --- |
| [项目总览](https://github.com/sihuangtech/probelune) | 项目总览、架构与部署指南 |
| [网页面板](https://github.com/sihuangtech/probelune-dashboard) | 网页控制台与 Go API 服务端 |
| [探测节点](https://github.com/sihuangtech/probelune-node) | 部署在监测网络中的探测执行器 |
| [移动客户端](https://github.com/sihuangtech/probelune-mobile) | Flutter Android 与 iOS 移动客户端 |
| [桌面客户端](https://github.com/sihuangtech/probelune-desktop) | Electron 桌面客户端与诊断工具 |
| [协议与数据定义](https://github.com/sihuangtech/probelune-contracts) | 共享 Go 数据结构、参数校验与接口协议 |


## 星标历史

[![星标历史](https://api.star-history.com/svg?repos=sihuangtech/probelune-contracts&type=Date)](https://star-history.com/#sihuangtech/probelune-contracts&Date)

## 联系开发者

- 开发者: [Snake Konginchrist](https://github.com/Snake-Konginchrist)
- 邮箱: [developer@skstudio.cn](mailto:developer@skstudio.cn)
- 问题反馈: [GitHub Issues](https://github.com/sihuangtech/probelune-contracts/issues)
