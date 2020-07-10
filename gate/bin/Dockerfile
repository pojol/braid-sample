
FROM alpine

COPY ./gateway_linux /tmp/app/

WORKDIR /tmp/app

ENV SERVICE_TAGS=braid,gateway
ENV SERVICE_14001_NAME=gateway


EXPOSE 14222

ENTRYPOINT ./gateway_linux $0 $@
