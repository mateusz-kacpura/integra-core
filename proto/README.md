# integra-proto

Protobuf definitions and generated Go code used by the current ledger-service configuration.

## Structure

- `accounting/v1/invoice.proto` defines the InvoiceService and invoice messages.
- `accounting/v1/ledger.proto` defines the LedgerService double-entry accounting RPC service.
- `iam/v1/user.proto` defines the IdentityService authentication and authorization RPC service.
- `user/v1/profile.proto` defines the ProfileService user profile RPC service.
- `gen/go/` contains generated Protobuf messages and Connect service code.
- `.github/workflows/` validates the schemas in CI.
- `.vscode/` configures the local Go development environment.

## Usage

```go
import (
    accountingv1 "github.com/mateusz-kacpura/integra-enterprise/integra-proto/gen/go/accounting/v1"
    "github.com/mateusz-kacpura/integra-enterprise/integra-proto/gen/go/accounting/v1/accountingv1connect"
)
```

Local development resolves modules through the root `go.work`. CI validates all nine modules with workspace disabled and disposable modfiles; permanent manifests contain no local `replace` directives.

## Commands

```bash
buf generate
buf lint
go build ./...
```

Generated code is committed so consumers do not need Protobuf generation tools.
Regenerate it after editing the schemas. This platform is a prototype with no external consumers: identity profile fields and deprecated getters have been removed and the remaining fields renumbered. CI checks schema lint and reproducible generation. Once contracts are released to external consumers, introduce a release compatibility baseline and preserve field numbers.

