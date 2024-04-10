package proto

//go:generate protoc -I=. --go_out=. entry.proto
//go:generate protoc -I=. --go_out=. migration.proto
//go:generate protoc -I=. --go_out=. observerAcl.proto
