# Test cases for checking migration

## TOC
- [Test cases for checking migration](#test-cases-for-checking-migration)
  - [TOC](#toc)
  - [Test cases](#test-cases)
    - [Preconditions](#preconditions)
    - [Migration tool](#migration-tool)
    - [Validation tool](#validation-tool)
    - [Observer](#observer)
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
  1. expand sandbox 1.4
  2. run test scripts/migration_test.go
    * save manually PublicKey, PrivateKey, Address in base58 format for all users
  3. do a migration
  4. expand new sandbox 2.4
  5. run scripts/migration_back_test.go

Expected result: all transactions completed successfully


* Verification migration on stand
  1. expand stand 1.4
  2. do more 200 transactions
  3. do a migration
  4. expand new stand 2.4

Expected result: hashes from 1.4 equal 2.4


* Verification migration on stand with script
  1. expand stand 1.4
  2. run test scripts/migration_test.go
    * save manually PublicKey, PrivateKey, Address in base58 format for all users
  3. do a migration
  4. expand new stand 2.4
  5. run scripts/migration_back_test.go

Expected result: all transactions completed successfully


### Validation tool

* Check Validation tool
  1. expand sandbox 1.4
  2. do some transactions
  3. run validation tool on hlf 1.4 (import transactions)
  4. do a migration 
  5. run validation tool on hlf 2.4 (export transactions)
  6. run validation tool for compare hashes from 1.4 and 2.4

Expected result: hashes from 1.4 equal 2.4


* Check comparison
  1. expand sandbox 1.4
  2. do some transactions
  3. run validation tool on hlf 1.4
  4. do a migration (export and import different transactions)
  5. run validation tool on hlf 2.4
  6. run validation tool for compare hashes from 1.4 and 2.4

Expected result: hashes from 1.4 not equal 2.4


### Observer

* Verification migration observer on stand
  1. expand stand 1.4
  2. take ch-dev base and integrate it into the stand 1.4
  3. do a migration
  4. expand new stand 2.4 with migration chain codes

Expected result: base on 2.4 contains hlf state from 1.4 (acl, user's balances)


## License

Apache-2.0

## Links

* [origin](https://github.com/anoideaopen/migrationcc)
