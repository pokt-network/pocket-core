FROM golang:1.17-alpine as build

WORKDIR /usr/src/app/pocket-core

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -o ./bin/pocket app/cmd/pocket_core/main.go && \
    cp ./bin/pocket /usr/local/bin/pocket

FROM alpine:latest

RUN apk add --no-cache ca-certificates

RUN mkdir -p /etc/pocket/
COPY --from=build /usr/local/bin/pocket /usr/local/bin

CMD ["pocket"]



