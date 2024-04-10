package utils

import (
	"crypto/sha256"
	"fmt"
)

const (
	ExportChunkKvEvent        = "exportChunkKV"
	ExportEndEvent            = "exportEnd"
	ImportChunkKvEvent        = "importChunkKV"
	CommitMigrationInfoEvent  = "commitMigrationInfo"
	InsertObserverACLRowEvent = "insertObserverACLRowEvent"
)

// ComputeHash returns the SHA256 checksum of the data.
func ComputeHash(data []byte) string {
	sum := sha256.Sum256(data)
	hash := fmt.Sprintf("%x", sum)

	fmt.Println("ComputeHash")
	fmt.Println(len(data))
	fmt.Println(hash)
	return hash
}
