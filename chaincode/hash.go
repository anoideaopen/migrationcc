package chaincode

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

const (
	ExportChunkKvEvent = "exportChunkKV"
	ImportChunkKvEvent = "importChunkKV"
)

// ComputeHash returns the SHA256 checksum of the data.
func ComputeHash(data []byte) string {
	sum := sha256.Sum256(data)
	hash := hex.EncodeToString(sum[:])

	fmt.Println("ComputeHash")
	fmt.Println(len(data))
	fmt.Println(hash)
	return hash
}
