# Probelune Contracts

[简体中文](README.zh-CN.md) · [Product and PRD](https://github.com/sihuangtech/probelune)

Versioned Go models, validation and network access policy for **Probelune (巡汐)**, a self-hosted multi-location monitoring system. This is a library and protocol specification, not a running service.

## Use

Requires Go 1.27 or newer. Integration settings are described in the [protocol documentation](docs/PROTOCOL.md).

```sh
export GONOPROXY=github.com/sihuangtech/probelune-contracts
export GONOSUMDB=github.com/sihuangtech/probelune-contracts
export GOMODCACHE="$PWD/.cache/go-mod-v0.1.1"
go get github.com/sihuangtech/probelune-contracts@v0.1.1
```

```go
check := contracts.Check{Type: "http", Target: "https://example.com", Timeout: 10}
if err := check.Validate(); err != nil {
    return err
}
```

`check.go` defines the check catalog and required capabilities; `models.go` defines monitors, expiring assignments, diagnostics and results. `validate.go` normalizes defaults and validates executor inputs. `network.go` validates all DNS answers before dialing a pinned IP; private/special networks require explicit CIDRs. Neither nodes nor the dashboard execute caller-supplied shell commands.

## Protocol

- [Node protocol and state semantics](docs/PROTOCOL.md)
- [Dashboard HTTP API](docs/API.md)
- [Check JSON Schema](schemas/check.schema.json)
- [OpenAPI entry point](openapi.yaml)

Checks include HTTP/HTTPS, TCP, UDP, ICMP, DNS and route diagnostics. The protocol carries node capabilities, configurable schedules, diagnostic runs and historical availability summaries. Applications consume the shared JSON wire representation while keeping their own user interfaces.

Protocol version `1` rejects incompatible nodes at polling time. Additive optional fields may be introduced within v1; breaking wire changes require a new protocol major. Module releases follow semantic versioning.

## Development

```sh
go fmt ./...
go vet ./...
go test ./...
```

Release notes are maintained in [GitHub Releases](https://github.com/sihuangtech/probelune-contracts/releases).

See [CONTRIBUTING.md](CONTRIBUTING.md) and [SECURITY.md](SECURITY.md). Licensed under GPL-3.0-or-later; see [LICENSE](LICENSE).


## Star History

[![Star history](https://api.star-history.com/svg?repos=sihuangtech/probelune-contracts&type=Date)](https://star-history.com/#sihuangtech/probelune-contracts&Date)

## Contact

- Developer: [Snake Konginchrist](https://github.com/Snake-Konginchrist)
- Email: [developer@skstudio.cn](mailto:developer@skstudio.cn)
- Feedback: [GitHub Issues](https://github.com/sihuangtech/probelune-contracts/issues)
