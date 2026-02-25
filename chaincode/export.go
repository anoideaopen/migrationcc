package chaincode

import (
	"errors"
	"fmt"

	"github.com/anoideaopen/migrationcc/proto"
	"github.com/hyperledger/fabric-chaincode-go/v2/shim"
	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
	"github.com/hyperledger/fabric-protos-go-apiv2/peer"
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

	hash := computeHash(bytes)
	fname := string(ctx.GetStub().GetArgs()[0])
	err = ctx.GetStub().SetEvent(fname, []byte(hash))
	if err != nil {
		return "", fmt.Errorf("couldn't set event: %w", err)
	}

	return string(bytes), nil
}

// getKVPairs get entries by page
// pageSize - keys between the bookmark (inclusive) and endKey (exclusive)
// bookmark - When an empty string is passed as a value to the bookmark argument, the returned
// iterator can be used to fetch the first `pageSize` keys between the startKey
// (inclusive) and endKey (exclusive).
// When the bookmark is a non-emptry string, the iterator can be used to fetch
// the first `pageSize` keys between the bookmark (inclusive) and endKey (exclusive).
// Note that only the bookmark present in a prior page of query results (ResponseMetadata)
// can be used as a value to the bookmark argument. Otherwise, an empty string must
// be passed as bookmark.
// isComposit - return simple keys or composit keys
func getKVPairs(stub shim.ChaincodeStubInterface, pageSize int32, bookmark string, isComposit bool) (*proto.Entries, error) {
	if pageSize < 0 {
		return nil, errors.New("pageSize need to be positive number")
	}

	var (
		iter     shim.StateQueryIteratorInterface
		metadata *peer.QueryResponseMetadata
		err      error
	)

	if isComposit {
		iter, metadata, err = stub.GetAllStatesCompositeKeyWithPagination(pageSize, bookmark)
	} else {
		iter, metadata, err = stub.GetStateByRangeWithPagination("", "", pageSize, bookmark)
	}
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = iter.Close()
	}()

	entries := new(proto.Entries)
	for iter.HasNext() {
		kv, err := iter.Next()
		if err != nil {
			return nil, err
		}

		entries.Entries = append(entries.Entries, &proto.Entry{
			Key:   kv.GetKey(),
			Value: kv.GetValue(),
		})
	}

	if len(entries.GetEntries()) == 0 {
		return nil, errors.New("not found kv pairs")
	}

	if metadata != nil {
		entries.Bookmark = metadata.GetBookmark()
	}

	return entries, nil
}
