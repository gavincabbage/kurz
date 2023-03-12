ARG GO_VERSION=1.20

FROM golang:${GO_VERSION}-alpine as build

RUN apk update
RUN apk add --no-cache build-base git

WORKDIR /

COPY go.mod .
#COPY go.sum . # No dependencies.

RUN go mod download

COPY . .

RUN go build -o kurzd .

FROM alpine

COPY --from=build kurzd .