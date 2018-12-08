FROM golang:1.10

ADD . /go/src/goophr/librarian

WORKDIR /go/src/goophr/librarian

RUN go install goophr/librarian

ENTRYPOINT /go/bin/librarian