@echo off

go build -o .bin/server sokoban/cmd/server
set GOOS=js
set GOARCH=wasm
go build -o cmd/server/resources/client.wasm sokoban/cmd/client