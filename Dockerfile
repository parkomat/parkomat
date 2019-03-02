FROM golang:1.10-alpine

RUN mkdir -p /go/src/github.com/parkomat/parkomat

COPY . /go/src/github.com/parkomat/parkomat

RUN apk add --update git

RUN mkdir -p /go/src/golang.org/x/

RUN git clone https://github.com/golang/sys.git /go/src/golang.org/x/sys

RUN git clone https://github.com/golang/text.git /go/src/golang.org/x/text

RUN git clone https://github.com/golang/net.git /go/src/golang.org/x/net
RUN go get github.com/parkomat/parkomat/...

EXPOSE 53
EXPOSE 53/udp
EXPOSE 80

CMD parkomat
