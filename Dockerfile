FROM golang:latest AS builder

ENV PROJECT_NAME="redirects"

RUN mkdir -p /go/src/$PROJECT_NAME
WORKDIR /go/src/$PROJECT_NAME

COPY ../../Desktop/f .

RUN ls
RUN GOOS=linux GOARCH=amd64 GO111MODULE=on CGO_ENABLED=0 go build -mod vendor -ldflags='-w -s' -o ./bin/binary cmd/redirect/main.go

FROM alpine:latest

ENV PROJECT_NAME="redirects"

COPY --from=builder /go/src/$PROJECT_NAME/bin/binary /binary
COPY --from=builder /go/src/$PROJECT_NAME/.env /.env

EXPOSE 5050

ENTRYPOINT ["/binary"]