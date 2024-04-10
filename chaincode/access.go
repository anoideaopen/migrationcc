package chaincode

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// KeyInitArgs chaincodes use this key to put initialization arguments
const KeyInitArgs = "__init"

type InitArgs struct {
	AdminSKI        []byte
	ValidatorsCount int
	Validators      []string
	AllowedMspId    string //nolint:revive,stylecheck
}

func CheckAccess(ctx contractapi.TransactionContextInterface) error {
	clientIdentity := ctx.GetClientIdentity()
	x509Certificate, err := clientIdentity.GetX509Certificate()
	if err != nil {
		return err
	}
	pk, ok := x509Certificate.PublicKey.(*ecdsa.PublicKey)
	if !ok {
		return errors.New("failed cast Certificate.PublicKey to ecdsa.PublicKey")
	}
	creatorSKI := sha256.Sum256(elliptic.Marshal(pk.Curve, pk.X, pk.Y))

	migrationRule, err := GetMigrationRule(ctx.GetStub())
	if err != nil {
		return err
	}
	if migrationRule == nil {
		return errors.New("migration rule not found")
	}

	mspID, err := clientIdentity.GetMSPID()
	if err != nil {
		return err
	}
	if migrationRule.MigrationMspID != mspID {
		return fmt.Errorf("constraint execute chaincode by migration rule. found msp id '%s' but expected '%s'", mspID, migrationRule.MigrationMspID)
	}

	if !bytes.Equal(creatorSKI[:], migrationRule.MigrationUserSKI) {
		return fmt.Errorf("constraint execute chaincode by migration rule. only specified certificate can invoke. found '%s' but expected '%s'", hex.EncodeToString(creatorSKI[:]), hex.EncodeToString(migrationRule.MigrationUserSKI))
	}

	return nil
}
