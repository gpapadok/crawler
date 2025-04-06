FROM golang:1.24

WORKDIR /go/src

COPY . .
RUN go build -o /go/bin/node cmd/node/node.go

CMD ["/go/bin/node"]
