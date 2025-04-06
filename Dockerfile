FROM golang:1.24-alpine3.21 AS build

WORKDIR /go/src

COPY . .
RUN go build -o /go/bin/node cmd/node/main.go

FROM alpine:3.21

COPY --from=build /go/bin/node /go/bin/node
COPY --from=build /go/src/.env /go/src/.env
WORKDIR /go/src

CMD ["/go/bin/node"]
