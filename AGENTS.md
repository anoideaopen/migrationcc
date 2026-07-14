# migrationcc

Hyperledger Fabric chaincode for state migration between Fabric networks.
Export/import world state with SHA-256 integrity hashes.

## Commands

```sh
go test -count 1 ./...          # unit tests (cache disabled)
go fmt ./...                    # formatting required by CI
go fix ./...                    # update code for newer Go API versions
golangci-lint run               # uses golangci-lint v2 (config: .golangci.yml)
go mod tidy && git diff --exit-code go.mod  # must be a no-op
```

Regenerate protobuf (from `proto/`):
```sh
go generate ./proto/
```
or directly:
```sh
protoc -I=. --go_out=paths=source_relative:. entry.proto
```

## Entrypoint

`main.go` bootstraps `contractapi.NewChaincode(&chaincode.MigrationContract{})` then `cc.Start()`.

## Architecture

| Directory | Purpose |
|---|---|
| `chaincode/` | Contract logic + tests + counterfeiter mocks |
| `proto/` | Protobuf schema (`Entry`, `Entries`) + generated Go |
| `doc/` | Test scenarios |

Two invocable methods on `MigrationContract`:

- **`ExportChunkKV(ctx, pageSize int32, bookmark string, isComposit bool)`** – paginated read from world state, returns serialized `Entries` + SHA-256 hash in chaincode event.
- **`ImportChunkKV(ctx, entriesStr string)`** – unmarshals protobuf, writes entries via `PutState`, emits hash event.

## Known gotchas

- **`FinishWriteBatch()` never called** in `ImportChunkKV` (`chaincode/import.go:14`). Batch is started but never committed – likely a bug.
- **`exportEnd` and `commitMigrationInfo`** are documented in README API but **not implemented** in code.
- **Event name** derived from `ctx.GetStub().GetArgs()[0]` (raw arg), not from contract API function name.
- **`isComposit`** param misspelled (should be `isComposite`). Exposed as-is in the contract API.
- **LICENSE is MIT** (check the file). README incorrectly says Apache 2.0.
- **`golangci-lint` runs only on non-test files** (`.golangci.yml`: `tests: false`).

## Testing

- Uses `counterfeiter`-generated mocks in `chaincode/mock/`.
- Test pattern: create `MigrationContract{}`, wrap with `contractapi.NewChaincode`, invoke via `cc.Invoke(mockStub)`.
- Regenerate mocks: `go generate ./chaincode/` (uses `//go:generate counterfeiter -generate`).
- Tests verify `SetEvent` was called but do **not** verify hash values.

## CI pipeline (sequential)

1. check-cyrillic-comments – blocks Cyrillic chars in source files
2. validate-go – `go mod tidy` + `go fmt`, fails if dirty
3. golangci-lint
4. go-test-unit – `go test -count 1 ./...`
