
FROM golang:alpine

COPY ./login_linux /tmp/app/

WORKDIR /tmp/app

ENV SERVICE_TAGS=braid,login
ENV SERVICE_14101_NAME=login


EXPOSE 14222

ENTRYPOINT ./login_linux $0 $@
