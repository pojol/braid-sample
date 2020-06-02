
FROM alpine

COPY ./base_linux /tmp/app/

WORKDIR /tmp/app

ENV SERVICE_TAGS=braid,base
ENV SERVICE_14201_NAME=base


EXPOSE 1201

ENTRYPOINT ./base_linux $0 $@
