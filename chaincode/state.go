package chaincode

import (
	"errors"

	"github.com/anoideaopen/migrationcc/proto"
	pb "github.com/golang/protobuf/proto" //nolint:staticcheck
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-protos-go/ledger/queryresult"
)

const (
	MigrationRuleKey = "___migration_rule"
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
func getKVPairs(stub shim.ChaincodeStubInterface, pageSize int32, bookmark string, onlyKeys bool) (*proto.Entries, error) {
	if pageSize < 0 {
		return nil, errors.New("pageSize need to be positive number")
	}

	iter, metadata, err := stub.GetAllKeyValueFromStateWithPagination(pageSize, bookmark)
	if err != nil {
		return nil, err
	}
	if iter == nil {
		return nil, errors.New("unexpected response from GetStateByRangeWithPagination. iter StateQueryIteratorInterface can't be nil")
	}

	migrationRuleKey, err := CreateMigrationKey()
	if err != nil {
		return nil, err
	}

	entries := new(proto.Entries)
	for iter.HasNext() {
		var kv *queryresult.KV
		kv, err = iter.Next()
		if err != nil {
			return nil, err
		}

		if kv.Key != migrationRuleKey {
			elems := &proto.Entry{
				Key: kv.Key,
			}

			if !onlyKeys {
				elems.Value = kv.Value
			}

			entries.Entries = append(entries.Entries, elems)
		}
	}

	if len(entries.Entries) == 0 {
		return nil, errors.New("not found kv pairs")
	}

	if metadata != nil {
		entries.Bookmark = metadata.Bookmark
	}

	return entries, nil
}

func GetMigrationRule(stub shim.ChaincodeStubInterface) (*proto.MigrationRule, error) {
	key, err := CreateMigrationKey()
	if err != nil {
		return nil, err
	}

	bytes, err := stub.GetState(key)
	if err != nil {
		return nil, err
	}
	if len(bytes) == 0 {
		return nil, errors.New("migration rule doesn't exists")
	}
	migrationRule := proto.MigrationRule{}
	err = pb.Unmarshal(bytes, &migrationRule)
	if err != nil {
		return nil, err
	}
	return &migrationRule, nil
}

func PutMigrationRule(stub shim.ChaincodeStubInterface, migrationRule *proto.MigrationRule) error {
	bytes, err := pb.Marshal(migrationRule)
	if err != nil {
		return err
	}

	key, err := CreateMigrationKey()
	if err != nil {
		return err
	}

	return stub.PutState(key, bytes)
}

func CreateMigrationKey() (string, error) {
	key, err := shim.CreateCompositeKey(MigrationRuleKey, []string{MigrationRuleKey})
	if err != nil {
		return "", err
	}
	return key, nil
}
