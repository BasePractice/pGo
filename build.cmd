@echo off

set PATH=E:\Programs\tinygo\bin;%PATH%

go build -o .bin/server ./cmd/server
set GOOS=js
set GOARCH=wasm
go build -o cmd/server/resources/client.wasm ./cmd/client