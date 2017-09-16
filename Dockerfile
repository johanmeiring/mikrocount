FROM alpine

RUN apk add --update ca-certificates

COPY mikrocount /mikrocount

ENTRYPOINT ["/mikrocount"]

