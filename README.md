# migrationcc

Chaincode for migration from old to new hlf

## TOC

- [migrationcc](#migrationcc)
  - [TOC](#toc)
  - [Description](#description)
  - [API](#api)
    - [Init chaincode](#init-chaincode)
    - [Invoke requests](#invoke-requests)
  - [Links](#links)
  - [License](#license)

## Description

Chaincode for unloading the state from the channel and ensuring data migration between hlf networks

## API

### Init chaincode

Initializes migration chaincode setup user information for migration process

**Args:**
```
-c '{"Args":[migrationMspID,migrationUserSKI]}'
```

**Args:**

	[0] migrationMspID     - msp id for user who responsible for migration process
	[1] migrationUserSKI       - SKI HLF user who responsible for migration process

### Invoke requests

- **exportEnd** - end of migration calculate the total hash of all keys and values during migration
  - NO Batch tx
  - **Signed by:** migration user with migrationMspID and migrationUserSKI
  - **Args:** [pageSize int32]
  - **Applied changes:** write hash of all keys and values to chaincode event

- **exportChunkKV** - get a set of keys and values and write hash to chaincode event
  - NO Batch tx
  - **Signed by:** migration user with migrationMspID and migrationUserSKI
  - **Args:** [pageSize int32, bookmark string, onlyKeys bool]
  - **Applied changes:** write hash of keys and values to chaincode event

- **importChunkKV** - add a set of keys and values to the state and write the hash to the chaincode event
  - NO Batch tx
  - **Signed by:** migration user with migrationMspID and migrationUserSKI
  - **Args:** [entries proto.Entries]
  - **Applied changes:** write keys and values to state, hash of keys and values write to chaincode event

- **commitMigrationInfo** - commit info about previous hlf version
  - NO Batch tx
  - **Signed by:** migration user with migrationMspID and migrationUserSKI
  - **Args:** [migrationInfo proto.MigrationInfo]
  - **Applied changes:** write info about previous hlf version to chaincode event

## Links



## License

Apache 2.0
