FROM golang:1.10

ADD . /go/src/goophr/concierge

WORKDIR /go/src/goophr/concierge

RUN go install goophr/concierge

ENTRYPOINT /go/bin/concierge