FROM golang:1.24

WORKDIR /go/src

COPY . .
RUN go build -o /go/bin/node cmd/node/main.go

CMD ["/go/bin/node"]
