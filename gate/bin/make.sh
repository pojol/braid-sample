#! /bin/sh

rm gateway
echo "build gateway ..."
go build -o gateway /Users/pojol/work/gohome/src/braid-game/gate/main.go

rm gateway_linux
echo "build gateway_linux ..."
GOOS=linux GOARCH=amd64 go build -o gateway_linux /Users/pojol/work/gohome/src/braid-game/gate/main.go
