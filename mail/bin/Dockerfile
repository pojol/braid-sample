FROM alpine

COPY ./mail_linux /tmp/app/

WORKDIR /tmp/app

ENV SERVICE_TAGS=braid,mail
ENV SERVICE_14301_NAME=mail


EXPOSE 14222

ENTRYPOINT ./mail_linux $0 $@
