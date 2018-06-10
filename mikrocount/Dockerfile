# build stage
FROM golang:1.10-alpine AS builder
RUN apk add --update ca-certificates
RUN apk add --no-cache git && go get -u github.com/golang/dep/...
ADD . /go/src/mikrocount
WORKDIR /go/src/mikrocount
RUN dep ensure
RUN go build -o mikrocount

# final stage
FROM alpine
WORKDIR /app
COPY --from=builder /go/src/mikrocount/mikrocount /app/mikrocount
ENTRYPOINT ["./mikrocount"]
