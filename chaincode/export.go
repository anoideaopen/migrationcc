package chaincode

import (
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
	pb "google.golang.org/protobuf/proto"
)

// ExportChunkKV get a set of keys and values and write hash to chaincode event
func (cc *MigrationContract) ExportChunkKV(ctx contractapi.TransactionContextInterface, pageSize int32, bookmark string, isComposit bool) (string, error) {
	entries, err := getKVPairs(ctx.GetStub(), pageSize, bookmark, isComposit)
	if err != nil {
		return "", err
	}

	bytes, err := pb.Marshal(entries)
	if err != nil {
		return "", fmt.Errorf("couldn't marshal protobuf: %w", err)
	}

	hash := ComputeHash(bytes)
	err = ctx.GetStub().SetEvent(ExportChunkKvEvent, []byte(hash))
	if err != nil {
		return "", fmt.Errorf("couldn't set event: %w", err)
	}

	return string(bytes), nil
}
