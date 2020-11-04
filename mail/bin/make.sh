#! /bin/sh

rm mail
echo "build mail ..."
go build -o mail /Users/pojol/work/gohome/src/braid-game/mail/main.go

rm mail_linux
echo "build mail_linux ..."
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o mail_linux /Users/pojol/work/gohome/src/braid-game/mail/main.go


# build
docker build -t braid-game/mail . --no-cache
