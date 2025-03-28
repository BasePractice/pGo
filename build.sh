#!/bin/sh

go build -o .bin/server ./cmd/server
GOOS=js GOARCH=wasm go build -o cmd/server/resources/client.wasm ./cmd/client