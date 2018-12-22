FROM golang:1.11.4-alpine3.8 AS builder
WORKDIR /go/src/mirror
COPY . .
RUN go build -v

FROM docker:18.06.1-ce
WORKDIR /usr/local/bin
COPY --from=builder /go/src/mirror/mirror /usr/bin/