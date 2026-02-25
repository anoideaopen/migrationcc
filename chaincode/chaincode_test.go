package chaincode

import (
	"testing"

	"github.com/anoideaopen/migrationcc/chaincode/mock"
	"github.com/anoideaopen/migrationcc/proto"
	"github.com/hyperledger/fabric-chaincode-go/v2/shim"
	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
	"github.com/hyperledger/fabric-protos-go-apiv2/ledger/queryresult"
	"github.com/hyperledger/fabric-protos-go-apiv2/peer"
	"github.com/stretchr/testify/assert"
	pb "google.golang.org/protobuf/proto"
)

//go:generate counterfeiter -generate
//counterfeiter:generate -o mock/chaincode_stub.go --fake-name ChaincodeStub . chaincodeStub
type chaincodeStub interface { //nolint:unused
	shim.ChaincodeStubInterface
}

//counterfeiter:generate -o mock/state_iterator.go --fake-name StateIterator . stateIterator
type stateIterator interface { //nolint:unused
	shim.StateQueryIteratorInterface
}

const (
	exportChunkKVMethodName = "exportChunkKV"
	importChunkKVMethodName = "importChunkKV"

	bookmark          = ""
	testPageSizeArg   = "20"
	testSimpleKeys    = "false"
	testCompositeKeys = "true"
)

// TestMigrationExportNonComposite test migration function export with mock stub
func TestMigrationExportNonComposite(t *testing.T) {
	migrationContract := MigrationContract{}
	var interfaces []contractapi.ContractInterface
	interfaces = append(interfaces, &migrationContract)

	exportCC, err := contractapi.NewChaincode(interfaces...)
	assert.NoError(t, err)
	mockStub := new(mock.ChaincodeStub)
	mockStub.GetFunctionAndParametersReturns(exportChunkKVMethodName, []string{testPageSizeArg, bookmark, testSimpleKeys})
	mockStub.GetArgsReturns([][]byte{[]byte(exportChunkKVMethodName), []byte(testPageSizeArg), []byte(bookmark), []byte(testSimpleKeys)})

	fakeIterator := &mock.StateIterator{}
	fakeIterator.HasNextReturnsOnCall(0, true)
	fakeIterator.HasNextReturnsOnCall(1, false)
	fakeIterator.NextReturns(&queryresult.KV{
		Key:   "key1",
		Value: []byte("value1"),
	}, nil)
	mockStub.GetStateByRangeWithPaginationReturns(fakeIterator, &peer.QueryResponseMetadata{
		FetchedRecordsCount: 1,
		Bookmark:            "",
	}, nil)

	resp := exportCC.Invoke(mockStub)

	entries := proto.Entries{}
	err = pb.Unmarshal(resp.GetPayload(), &entries)
	assert.NoError(t, err)
	pb.Equal(&entries, &proto.Entries{
		Entries: []*proto.Entry{
			{Key: "key1", Value: []byte("value1")},
		},
	})
	assert.Equal(t, 1, mockStub.SetEventCallCount())
	evName, _ := mockStub.SetEventArgsForCall(0)
	assert.Equal(t, exportChunkKVMethodName, evName)
}

// TestMigrationExportComposite test migration function export with mock stub
func TestMigrationExportComposite(t *testing.T) {
	migrationContract := MigrationContract{}
	var interfaces []contractapi.ContractInterface
	interfaces = append(interfaces, &migrationContract)

	exportCC, err := contractapi.NewChaincode(interfaces...)
	assert.NoError(t, err)
	mockStub := new(mock.ChaincodeStub)
	mockStub.GetFunctionAndParametersReturns(exportChunkKVMethodName, []string{testPageSizeArg, bookmark, testCompositeKeys})
	mockStub.GetArgsReturns([][]byte{[]byte(exportChunkKVMethodName), []byte(testPageSizeArg), []byte(bookmark), []byte(testCompositeKeys)})

	key, err := shim.CreateCompositeKey("type1", []string{"key1"})
	assert.NoError(t, err)
	fakeIterator := &mock.StateIterator{}
	fakeIterator.HasNextReturnsOnCall(0, true)
	fakeIterator.HasNextReturnsOnCall(1, false)
	fakeIterator.NextReturns(&queryresult.KV{
		Key:   key,
		Value: []byte("value1"),
	}, nil)
	mockStub.GetAllStatesCompositeKeyWithPaginationReturns(fakeIterator, &peer.QueryResponseMetadata{
		FetchedRecordsCount: 1,
		Bookmark:            "",
	}, nil)

	resp := exportCC.Invoke(mockStub)

	entries := proto.Entries{}
	err = pb.Unmarshal(resp.GetPayload(), &entries)
	assert.NoError(t, err)
	pb.Equal(&entries, &proto.Entries{
		Entries: []*proto.Entry{
			{Key: key, Value: []byte("value1")},
		},
	})
	assert.Equal(t, 1, mockStub.SetEventCallCount())
	evName, _ := mockStub.SetEventArgsForCall(0)
	assert.Equal(t, exportChunkKVMethodName, evName)
}

// TestMigrationImport test migration function import with mock stub
func TestMigrationImport(t *testing.T) {
	migrationContract := MigrationContract{}
	var interfaces []contractapi.ContractInterface
	interfaces = append(interfaces, &migrationContract)

	importCC, err := contractapi.NewChaincode(interfaces...)
	assert.NoError(t, err)
	mockStub := new(mock.ChaincodeStub)
	key, err := shim.CreateCompositeKey("type1", []string{"key2"})
	assert.NoError(t, err)
	entries := &proto.Entries{
		Entries: []*proto.Entry{
			{Key: "key1", Value: []byte("value1")},
			{Key: key, Value: []byte("value2")},
		},
	}
	entriesBites, err := pb.Marshal(entries)
	assert.NoError(t, err)
	mockStub.GetFunctionAndParametersReturns(importChunkKVMethodName, []string{string(entriesBites)})
	mockStub.GetArgsReturns([][]byte{[]byte(importChunkKVMethodName), entriesBites})

	resp := importCC.Invoke(mockStub)
	assert.Len(t, resp.GetPayload(), 0)
	assert.Equal(t, 2, mockStub.PutStateCallCount())
	k, v := mockStub.PutStateArgsForCall(0)
	assert.Equal(t, "key1", k)
	assert.Equal(t, []byte("value1"), v)
	k, v = mockStub.PutStateArgsForCall(1)
	assert.Equal(t, key, k)
	assert.Equal(t, []byte("value2"), v)
	assert.Equal(t, 1, mockStub.SetEventCallCount())
	evName, _ := mockStub.SetEventArgsForCall(0)
	assert.Equal(t, importChunkKVMethodName, evName)
}
