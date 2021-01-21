FROM golang:alpine as build-env

ENV GO111MODULE=on

RUN apk update && apk add bash ca-certificates git gcc g++ libc-dev

RUN mkdir /message_service

WORKDIR /message_service

COPY . .

RUN go build -o message_service .

CMD ./message_service