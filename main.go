package main

import (
	"github.com/anoideaopen/migrationcc/chaincode"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

func main() {
	migrationContract := new(chaincode.MigrationContract)
	cc, err := contractapi.NewChaincode(migrationContract)
	if err != nil {
		panic(err.Error())
	}
	if err = cc.Start(); err != nil {
		panic(err.Error())
	}
}
