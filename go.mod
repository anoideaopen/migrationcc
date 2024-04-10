module github.com/anoideaopen/migrationcc

go 1.14

require (
	github.com/golang/protobuf v1.5.2
	github.com/google/uuid v1.3.0
	github.com/hyperledger/fabric-chaincode-go v0.0.0-20220720122508-9207360bbddd
	github.com/hyperledger/fabric-contract-api-go v1.1.1
	github.com/hyperledger/fabric-protos-go v0.0.0-20220827195505-ce4c067a561d
	github.com/stretchr/testify v1.8.0
	google.golang.org/protobuf v1.28.1
)

replace github.com/hyperledger/fabric-chaincode-go => github.com/scientificideas/fabric-chaincode-go v0.0.0-20221010114645-63470279a8e2
