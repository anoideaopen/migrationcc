package chaincode

import (
	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
)

type MigrationContract struct {
	contractapi.Contract
}

func (cc *MigrationContract) Init(ctx contractapi.TransactionContextInterface) error {
	return nil
}
