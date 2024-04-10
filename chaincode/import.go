package chaincode

import (
	"errors"
	"fmt"

	"github.com/anoideaopen/migrationcc/proto"
	"github.com/anoideaopen/migrationcc/utils"
	pb "github.com/golang/protobuf/proto" //nolint:staticcheck
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// ImportChunkKV add a set of keys and values to the state and write the hash to the chaincode event
func (cc *MigrationContract) ImportChunkKV(ctx contractapi.TransactionContextInterface, entriesStr string) error {
	err := CheckAccess(ctx)
	if err != nil {
		return err
	}
	if len(entriesStr) == 0 {
		return errors.New("arg entries can't be empty")
	}

	entriesBytes := []byte(entriesStr)
	entries := proto.Entries{}
	err = pb.Unmarshal(entriesBytes, &entries)
	if err != nil {
		return err
	}

	if len(entries.Entries) == 0 {
		return errors.New("wrong count of Entries, entries is empty")
	}

	for _, entry := range entries.Entries {
		err = ctx.GetStub().PutState(entry.Key, entry.Value)
		if err != nil {
			return err
		}
	}

	hash := utils.ComputeHash(entriesBytes)
	err = ctx.GetStub().SetEvent(utils.ImportChunkKvEvent, []byte(hash))
	if err != nil {
		return fmt.Errorf("couldn't set event: %w", err)
	}

	return nil
}

// CommitMigrationInfo commit info about previous hlf version
func (cc *MigrationContract) CommitMigrationInfo(ctx contractapi.TransactionContextInterface, migrationInfoStr string) error {
	err := CheckAccess(ctx)
	if err != nil {
		return err
	}

	if len(migrationInfoStr) == 0 {
		return errors.New("arg migration info can't be empty")
	}

	migrationInfoBytes := []byte(migrationInfoStr)
	migrationInfo := &proto.MigrationInfo{}
	err = pb.Unmarshal(migrationInfoBytes, migrationInfo)
	if err != nil {
		return err
	}
	if len(migrationInfo.BlockHash) == 0 {
		return errors.New("migrationInfo.BlockHash can't be empty")
	}
	if migrationInfo.BlockNumber == 0 {
		return errors.New("migrationInfo.BlockNumber can't be zero")
	}
	if len(migrationInfo.BlockDataHash) == 0 {
		return errors.New("migrationInfo.BlockDataHash can't be empty")
	}
	return ctx.GetStub().SetEvent(utils.CommitMigrationInfoEvent, migrationInfoBytes)
}
