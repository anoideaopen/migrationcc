package chaincode

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"
	"testing"

	"github.com/anoideaopen/migrationcc/mock"
	proto2 "github.com/anoideaopen/migrationcc/proto"
	"github.com/anoideaopen/migrationcc/utils"
	pb "github.com/golang/protobuf/proto" //nolint:staticcheck
	"github.com/google/uuid"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/hyperledger/fabric-protos-go/msp"
	"github.com/hyperledger/fabric-protos-go/peer"
	"github.com/stretchr/testify/assert"
)

const (
	adminMSP              = "adminMSP"
	adminCert             = "-----BEGIN CERTIFICATE-----\nMIICSDCCAe6gAwIBAgIQAJwYy5PJAYSC1i0UgVN5bjAKBggqhkjOPQQDAjCBhzEL\nMAkGA1UEBhMCVVMxEzARBgNVBAgTCkNhbGlmb3JuaWExFjAUBgNVBAcTDVNhbiBG\ncmFuY2lzY28xIzAhBgNVBAoTGmF0b215emUudWF0LmRsdC5hdG9teXplLmNoMSYw\nJAYDVQQDEx1jYS5hdG9teXplLnVhdC5kbHQuYXRvbXl6ZS5jaDAeFw0yMDEwMTMw\nODU2MDBaFw0zMDEwMTEwODU2MDBaMHUxCzAJBgNVBAYTAlVTMRMwEQYDVQQIEwpD\nYWxpZm9ybmlhMRYwFAYDVQQHEw1TYW4gRnJhbmNpc2NvMQ4wDAYDVQQLEwVhZG1p\nbjEpMCcGA1UEAwwgQWRtaW5AYXRvbXl6ZS51YXQuZGx0LmF0b215emUuY2gwWTAT\nBgcqhkjOPQIBBggqhkjOPQMBBwNCAAQGQX9IhgjCtd3mYZ9DUszmUgvubepVMPD5\nFlwjCglB2SiWuE2rT/T5tHJsU/Y9ZXFtOOpy/g9tQ/0wxDWwpkbro00wSzAOBgNV\nHQ8BAf8EBAMCB4AwDAYDVR0TAQH/BAIwADArBgNVHSMEJDAigCBSv0ueZaB3qWu/\nAwOtbOjaLd68woAqAklfKKhfu10K+DAKBggqhkjOPQQDAgNIADBFAiEAoKRQLe4U\nFfAAwQs3RCWpevOPq+J8T4KEsYvswKjzfJYCIAs2kOmN/AsVUF63unXJY0k9ktfD\nfAaqNRaboY1Yg1iQ\n-----END CERTIFICATE-----"
	importChaincodeName   = "importcc"
	exportChaincodeName   = "exportcc"
	bookmark              = "bookmark"
	countStateDBRecords   = 2000
	successResponseStatus = int32(200)

	exportChunkKVMethodName       = "exportChunkKV"
	exportEndMethodName           = "exportEnd"
	commitMigrationInfoMethodName = "commitMigrationInfo"
	importChunkKVMethodName       = "importChunkKV"

	testPageSizeArg    = "20"
	testOnlyKeysArg    = "true"
	testNotOnlyKeysArg = "false"
)

// TestMigrationImportExport test migration function export and import with mock stub
func TestMigrationImportExport(t *testing.T) {
	cert := parseTestCert(t)

	mockStubExportCC := createMigrationChaincode(t, exportChaincodeName, cert)

	countKVPair := 2000
	putTestData(t, mockStubExportCC, countKVPair)

	resp := invokeExportChunkKV(t, mockStubExportCC)

	invokeExportChunkKVOnlyKeys(t, mockStubExportCC)

	invokeExportEnd(t, mockStubExportCC)

	mockStubImportCC := createMigrationChaincode(t, importChaincodeName, cert)

	assert.NotEqual(t, mockStubExportCC.State, mockStubImportCC.State)

	invokeImportChunkKV(t, mockStubImportCC, resp)

	assert.Equal(t, mockStubExportCC.State, mockStubImportCC.State)
}

// TestExportEnd export check common state from all data
func TestExportEnd(t *testing.T) {
	cert := parseTestCert(t)

	mockStubExportCC := createMigrationChaincode(t, exportChaincodeName, cert)

	countKVPair := 2000
	putTestData(t, mockStubExportCC, countKVPair)

	resp := invokeExportChunkKV(t, mockStubExportCC)

	var allChunkHash string
	allChunkHash += utils.ComputeHash(resp.Payload)
	hash := utils.ComputeHash([]byte(allChunkHash))

	res1ExportEnd := invokeExportEnd(t, mockStubExportCC)
	res2ExportEnd := invokeExportEnd(t, mockStubExportCC)
	if res1ExportEnd != res2ExportEnd {
		assert.Fail(t, "method exportEnd give different result")
	}
	if res1ExportEnd != hash {
		assert.Fail(t, "common hash exportEnd can't be reproduced from source response exportChunkKV")
	}

	var eventExportChunkKV *peer.ChaincodeEvent
	var eventExportEnd1 *peer.ChaincodeEvent
	var eventExportEnd2 *peer.ChaincodeEvent
	found := len(mockStubExportCC.ChaincodeEventsChannel)
	expected := 3
	if found == expected {
		eventExportChunkKV = <-mockStubExportCC.ChaincodeEventsChannel
		eventExportEnd1 = <-mockStubExportCC.ChaincodeEventsChannel
		eventExportEnd2 = <-mockStubExportCC.ChaincodeEventsChannel
	} else {
		assert.Fail(t, fmt.Sprintf("problem to recive needed events. found %d but expected %d", found, expected))
	}

	if string(eventExportEnd1.Payload) != string(eventExportEnd2.Payload) {
		assert.Fail(t, "event after repeat 'exportEnd' not equals")
	}
	if string(eventExportEnd1.Payload) != utils.ComputeHash(eventExportChunkKV.Payload) {
		assert.Fail(t, "event ")
	}
}

// TestMigrationInfo test put information about last block from previous ledger
func TestMigrationInfo(t *testing.T) {
	cert := parseTestCert(t)

	mockStubImportCC := createMigrationChaincode(t, importChaincodeName, cert)

	migrationInfo := &proto2.MigrationInfo{
		BlockNumber:   1,
		BlockHash:     []byte("2"),
		BlockDataHash: []byte("3"),
	}
	invokeCommitMigrationInfo(t, mockStubImportCC, migrationInfo)
}

func invokeImportChunkKV(t *testing.T, mockStub *mock.Stub, resp peer.Response) {
	mockStub.TxID = txIDGen()
	args := [][]byte{[]byte(importChunkKVMethodName), resp.Payload}
	resp = mockStub.InvokeChaincode(mockStub.Name, args, mockStub.Name)
	assert.NotNil(t, resp)
	assert.Equal(t, resp.Status, successResponseStatus)
}

func invokeExportEnd(t *testing.T, mockStub *mock.Stub) string {
	mockStub.TxID = txIDGen()
	args := [][]byte{[]byte(exportEndMethodName), []byte(testPageSizeArg)}
	resp := mockStub.InvokeChaincode(mockStub.Name, args, mockStub.Name)
	assert.NotNil(t, resp)
	assert.Equal(t, resp.Status, successResponseStatus)
	return string(resp.Payload)
}

func invokeCommitMigrationInfo(t *testing.T, mockStub *mock.Stub, migrationInfo *proto2.MigrationInfo) {
	mockStub.TxID = txIDGen()
	bytes, err := pb.Marshal(migrationInfo)
	assert.NoError(t, err)
	args := [][]byte{[]byte(commitMigrationInfoMethodName), bytes}
	resp := mockStub.InvokeChaincode(mockStub.Name, args, mockStub.Name)
	assert.NotNil(t, resp)
	assert.Equal(t, resp.Status, successResponseStatus)
}

func invokeExportChunkKV(t *testing.T, mockStubExportcc *mock.Stub) peer.Response {
	mockStubExportcc.TxID = txIDGen()
	args := [][]byte{[]byte(exportChunkKVMethodName), []byte(testPageSizeArg), []byte(bookmark), []byte(testNotOnlyKeysArg)}
	resp := mockStubExportcc.InvokeChaincode(exportChaincodeName, args, exportChaincodeName)
	assert.NotNil(t, resp)
	assert.Equal(t, resp.Status, successResponseStatus)
	assert.NotEmpty(t, resp.Payload)

	entries := proto2.Entries{}
	err := pb.Unmarshal(resp.Payload, &entries)
	assert.NoError(t, err)
	assert.Equal(t, countStateDBRecords, len(entries.Entries))

	for _, entry := range entries.Entries {
		assert.NotEmpty(t, entry.Key)
		assert.NotEmpty(t, entry.Value)
	}
	return resp
}

func invokeExportChunkKVOnlyKeys(t *testing.T, mockStubExportcc *mock.Stub) peer.Response {
	mockStubExportcc.TxID = txIDGen()
	args := [][]byte{[]byte(exportChunkKVMethodName), []byte(testPageSizeArg), []byte(bookmark), []byte(testOnlyKeysArg)}
	resp := mockStubExportcc.InvokeChaincode(exportChaincodeName, args, exportChaincodeName)
	assert.NotNil(t, resp)
	assert.Equal(t, resp.Status, successResponseStatus)
	assert.NotEmpty(t, resp.Payload)

	entries := proto2.Entries{}
	err := pb.Unmarshal(resp.Payload, &entries)
	assert.NoError(t, err)
	assert.Equal(t, countStateDBRecords, len(entries.Entries))

	for _, entry := range entries.Entries {
		assert.NotEmpty(t, entry.Key)
		assert.Empty(t, entry.Value)
	}
	return resp
}

func createMigrationChaincode(t *testing.T, name string, cert *x509.Certificate) *mock.Stub {
	migrationContract := MigrationContract{}
	var interfaces []contractapi.ContractInterface
	interfaces = append(interfaces, &migrationContract)
	exportCC, err := contractapi.NewChaincode(interfaces...)
	assert.NoError(t, err)
	mockStubExportCC := mock.NewMockStub(name, exportCC)
	mockStubExportCC.MockPeerChaincode(name, mockStubExportCC, name)
	pk, ok := cert.PublicKey.(*ecdsa.PublicKey)
	if !ok {
		assert.Fail(t, "failed cast Certificate.PublicKey to ecdsa.PublicKey")
	}
	creatorSKI := sha256.Sum256(elliptic.Marshal(pk.Curve, pk.X, pk.Y))
	mockStubExportCC.MockInit(txIDGen(), [][]byte{[]byte("init"), []byte(adminMSP), []byte(hex.EncodeToString(creatorSKI[:]))})
	err = SetCreator(mockStubExportCC, adminMSP, cert.Raw)
	assert.NoError(t, err)
	return mockStubExportCC
}

func parseTestCert(t *testing.T) *x509.Certificate {
	pcert, _ := pem.Decode([]byte(adminCert))
	cert, err := x509.ParseCertificate(pcert.Bytes)
	assert.NoError(t, err)
	return cert
}

func putTestData(t *testing.T, mockStubExportCC *mock.Stub, countKVPair int) {
	for i := 0; i < countKVPair; i++ {
		str := fmt.Sprint(i)
		mockStubExportCC.TxID = txIDGen()
		err := mockStubExportCC.PutState(str, []byte(str))
		assert.NoError(t, err)
	}
}

func txIDGen() string {
	txID := [16]byte(uuid.New())
	return hex.EncodeToString(txID[:])
}

func marshalIdentity(creatorMSP string, creatorCert []byte) ([]byte, error) {
	pemblock := &pem.Block{Type: "CERTIFICATE", Bytes: creatorCert}
	pemBytes := pem.EncodeToMemory(pemblock)
	if pemBytes == nil {
		return nil, errors.New("encoding of identity failed")
	}

	creator := &msp.SerializedIdentity{Mspid: creatorMSP, IdBytes: pemBytes}
	marshaledIdentity, err := pb.Marshal(creator)
	if err != nil {
		return nil, err
	}
	return marshaledIdentity, nil
}

func SetCreator(stub *mock.Stub, creatorMSP string, creatorCert []byte) error {
	marshaledIdentity, err := marshalIdentity(creatorMSP, creatorCert)
	if err != nil {
		return err
	}
	stub.Creator = marshaledIdentity
	return nil
}
