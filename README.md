# Probelune Contracts

[简体中文](README.zh-CN.md) · [Product and PRD](https://github.com/sihuangtech/probelune)

Versioned Go models, validation and network access policy for **Probelune (巡汐)**, a self-hosted multi-location monitoring system. This is a library and protocol specification, not a running service.

## Use

Requires Go 1.27 or newer.

```sh
go get github.com/sihuangtech/probelune-contracts@v0.1.1
```

```go
check := contracts.Check{Type: "http", Target: "https://example.com", Timeout: 10}
if err := check.Validate(); err != nil {
    return err
}
```

`models.go` defines checks, monitors, expiring assignments, diagnostics and results. `validate.go` normalizes defaults and validates executor inputs. `network.go` validates all DNS answers before dialing a pinned IP; private/special networks require explicit CIDRs. Neither nodes nor the dashboard execute caller-supplied shell commands.

## Protocol

- [Node protocol and state semantics](docs/PROTOCOL.md)
- [Dashboard HTTP API](docs/API.md)
- [Check JSON Schema](schemas/check.schema.json)
- [OpenAPI entry point](openapi.yaml)

Go consumers import the tagged module. Vue, Electron and Flutter consume its JSON wire representation; they have independent application models and releases. There is no shared JavaScript workspace or sibling-directory dependency in published builds.

Protocol version `1` rejects incompatible nodes at polling time. Additive optional fields may be introduced within v1; breaking wire changes require a new protocol major. Module releases follow semantic versioning and retain immutable tags.

## Development

```sh
go fmt ./...
go vet ./...
go test ./...
```

See [CONTRIBUTING.md](CONTRIBUTING.md) and [SECURITY.md](SECURITY.md). Licensed under GPL-3.0-or-later; see [LICENSE](LICENSE).
