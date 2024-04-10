package chaincode

import (
	"fmt"

	"github.com/anoideaopen/migrationcc/utils"
	pb "github.com/golang/protobuf/proto" //nolint:staticcheck
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// ExportEnd end of migration calculate the total hash of all keys and values during migration
func (cc *MigrationContract) ExportEnd(ctx contractapi.TransactionContextInterface, pageSize int32) (string, error) {
	err := CheckAccess(ctx)
	if err != nil {
		return "", err
	}

	var (
		allChunkHash string
		bookmark     string
	)
	for {
		entries, err := getKVPairs(ctx.GetStub(), pageSize, bookmark, false)
		if err != nil {
			return "", err
		}

		bytes, err := pb.Marshal(entries)
		if err != nil {
			return "", fmt.Errorf("couldn't marshal protobuf: %w", err)
		}

		allChunkHash += utils.ComputeHash(bytes)

		if entries.Bookmark == "" {
			break
		}

		bookmark = entries.Bookmark
	}

	hash := utils.ComputeHash([]byte(allChunkHash))
	err = ctx.GetStub().SetEvent(utils.ExportEndEvent, []byte(hash))
	if err != nil {
		return "", fmt.Errorf("couldn't set event: %w", err)
	}

	return hash, nil
}

// ExportChunkKV get a set of keys and values and write hash to chaincode event
func (cc *MigrationContract) ExportChunkKV(ctx contractapi.TransactionContextInterface, pageSize int32, bookmark string, onlyKeys bool) (string, error) {
	err := CheckAccess(ctx)
	if err != nil {
		return "", err
	}

	entries, err := getKVPairs(ctx.GetStub(), pageSize, bookmark, onlyKeys)
	if err != nil {
		return "", err
	}

	bytes, err := pb.Marshal(entries)
	if err != nil {
		return "", fmt.Errorf("couldn't marshal protobuf: %w", err)
	}

	hash := utils.ComputeHash(bytes)
	err = ctx.GetStub().SetEvent(utils.ExportChunkKvEvent, []byte(hash))
	if err != nil {
		return "", fmt.Errorf("couldn't set event: %w", err)
	}

	return string(bytes), nil
}
