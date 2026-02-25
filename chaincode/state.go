package chaincode

import (
	"errors"

	"github.com/anoideaopen/migrationcc/proto"
	"github.com/hyperledger/fabric-chaincode-go/v2/shim"
	"github.com/hyperledger/fabric-protos-go-apiv2/peer"
)

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
// onlyKeys - return only key, without value
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
