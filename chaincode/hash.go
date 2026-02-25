package chaincode

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

// computeHash returns the SHA256 checksum of the data.
func computeHash(data []byte) string {
	sum := sha256.Sum256(data)
	hash := hex.EncodeToString(sum[:])

	fmt.Println("computeHash: ", hash)
	return hash
}
