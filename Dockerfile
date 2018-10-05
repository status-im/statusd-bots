FROM golang:1.10-alpine as builder

ARG PROJECT_PATH=/go/src/github.com/status-im/statusd-bots

RUN apk add --no-cache gcc musl-dev linux-headers make

RUN mkdir -p $PROJECT_PATH
ADD . $PROJECT_PATH
RUN cd $PROJECT_PATH && make build

FROM alpine:latest

ARG PROJECT_PATH=/go/src/github.com/status-im/statusd-bots

RUN apk add --no-cache ca-certificates bash

COPY --from=builder $PROJECT_PATH/bin/* /usr/local/bin/
