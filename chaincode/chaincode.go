package chaincode

import (
	"encoding/hex"
	"errors"

	"github.com/anoideaopen/migrationcc/proto"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type MigrationContract struct {
	contractapi.Contract
}

func (cc *MigrationContract) Init(ctx contractapi.TransactionContextInterface, migrationMspID string, migrationUserSKI string) error {
	if len(migrationMspID) == 0 {
		return errors.New("need to setup msp id for user who responsible for migration process")
	}
	if len(migrationUserSKI) == 0 {
		return errors.New("need to setup ski for user who responsible for migration process")
	}
	skiBytes, err := hex.DecodeString(migrationUserSKI)
	if err != nil {
		return err
	}
	migrationRule := &proto.MigrationRule{}
	migrationRule.MigrationUserSKI = skiBytes
	migrationRule.MigrationMspID = migrationMspID
	return PutMigrationRule(ctx.GetStub(), migrationRule)
}
