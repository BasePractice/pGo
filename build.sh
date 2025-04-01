#!/bin/sh

go build -o .bin/server ./cmd/server
go build -o .bin/desktop ./cmd/desktop
GOOS=js GOARCH=wasm go build -o cmd/server/resources/client.wasm ./cmd/client