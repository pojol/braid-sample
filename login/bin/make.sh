#! /bin/sh

rm login
echo "build login ..."
go build -race -o login /Users/pojol/work/gohome/src/braid-game/login/main.go

rm login_linux
echo "build login_linux ..."
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o login_linux /Users/pojol/work/gohome/src/braid-game/login/main.go


# build
docker build -t braid-game/login .
