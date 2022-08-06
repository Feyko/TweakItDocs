//go:build generate

package main

//go:generate protoc --go_out=./internal/translate/ internal/translate/proto.proto
