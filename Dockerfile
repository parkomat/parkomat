FROM golang:1.6-alpine

RUN mkdir -p /go/src/github.com/parkomat/parkomat

COPY . /go/src/github.com/parkomat/parkomat

RUN apk add --update git

RUN go get github.com/parkomat/parkomat/...

RUN go build github.com/parkomat/parkomat

EXPOSE 53
EXPOSE 53/udp

CMD parkomat
