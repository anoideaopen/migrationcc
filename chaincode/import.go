package chaincode

import (
	"errors"
	"fmt"

	"github.com/anoideaopen/migrationcc/proto"
	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
	pb "google.golang.org/protobuf/proto"
)

// ImportChunkKV add a set of keys and values to the state and write the hash to the chaincode event
func (cc *MigrationContract) ImportChunkKV(ctx contractapi.TransactionContextInterface, entriesStr string) error {
	if len(entriesStr) == 0 {
		return errors.New("arg entries can't be empty")
	}

	entriesBytes := []byte(entriesStr)
	entries := proto.Entries{}
	err := pb.Unmarshal(entriesBytes, &entries)
	if err != nil {
		return err
	}

	if len(entries.GetEntries()) == 0 {
		return errors.New("wrong count of Entries, entries is empty")
	}

	for _, entry := range entries.GetEntries() {
		err = ctx.GetStub().PutState(entry.GetKey(), entry.GetValue())
		if err != nil {
			return err
		}
	}

	hash := ComputeHash(entriesBytes)
	err = ctx.GetStub().SetEvent(ImportChunkKvEvent, []byte(hash))
	if err != nil {
		return fmt.Errorf("couldn't set event: %w", err)
	}

	return nil
}
