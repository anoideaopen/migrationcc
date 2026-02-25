# Test cases for checking migration

## TOC
- [Test cases for checking migration](#test-cases-for-checking-migration)
  - [TOC](#toc)
  - [Test cases](#test-cases)
    - [Preconditions](#preconditions)
    - [Migration tool](#migration-tool)
  - [License](#license)
  - [Links](#links)

## Test cases

### Preconditions

* Set env values for scripts:
  HLF_PROXY_URL (ex: HLF_PROXY_URL  = "http://localhost:9091")
  HLF_PROXY_AUTH_TOKEN (ex: HLF_PROXY_AUTH_TOKEN = "test")
  FIAT_ISSUER_PRIVATE_KEY

* Set values in scripts/utils/utils.go for scripts:
	TokenName (ex: "st" for stand and "fiat" for sandbox)
	TokenNAME (ex: "ST" for stand and "FIAT" for sandbox)
	SwapTokenName (ex: "ba" for stand and "cc" for sandbox)
	SwapTokenNAME (ex: "BA" for stand and "CC" for sandbox)

### Migration tool

* Verification migration on sandbox
  1. expand fabric old
  2. run test scripts/migration_test.go
    * save manually PublicKey, PrivateKey, Address in base58 format for all users
  3. do a migration
  4. expand fabric new
  5. run scripts/migration_back_test.go

Expected result: all transactions completed successfully


## License

Apache-2.0

## Links

* [origin](https://github.com/anoideaopen/migrationcc)
