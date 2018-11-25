FROM golang:1.11 AS build

WORKDIR /go/src/github.com/pokt-network/pocket-core/
RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

COPY . .
RUN make -B

FROM alpine:latest

RUN apk --no-cache add ca-certificates
WORKDIR /app/

COPY --from=builder /go/src/github.com/pokt-network/pocket-core/core .
CMD ["./core"]
