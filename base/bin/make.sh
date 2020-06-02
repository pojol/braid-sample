#! /bin/sh

rm base
echo "build base ..."
go build -o base /Users/pojol/work/gohome/src/braid-game/base/main.go

rm base_linux
echo "build base_linux ..."
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o base_linux /Users/pojol/work/gohome/src/braid-game/base/main.go
