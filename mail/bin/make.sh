#! /bin/sh

rm mail
echo "build mail ..."
go build -o mail /Users/pojol/work/gohome/src/braid-game/mail/main.go

rm mail_linux
echo "build mail_linux ..."
GOOS=linux GOARCH=amd64 go build -o mail_linux /Users/pojol/work/gohome/src/braid-game/mail/main.go
