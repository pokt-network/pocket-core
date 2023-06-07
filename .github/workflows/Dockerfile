# Based on a previous implementation to make sure we don't break existing deployments.
# https://github.com/pokt-network/pocket-core-deployments/blob/staging/docker/Dockerfile

FROM golang:1.17-alpine as build
RUN apk add --no-cache ca-certificates
WORKDIR /build
ADD . .
RUN go build -o pocket app/cmd/pocket_core/main.go

FROM alpine
RUN apk add --update --no-cache expect bash leveldb-dev tzdata && cp /usr/share/zoneinfo/America/New_York /etc/localtime \
    && addgroup --gid 1001 -S app \
    && adduser --uid 1005 -S -G app app
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=build /build/pocket /bin/pocket
COPY .github/workflows/entrypoint.sh /home/app/entrypoint.sh
RUN chown -R app /bin/pocket && mkdir -p /home/app/.pocket/config && chown -R app /home/app/.pocket
ENTRYPOINT ["/usr/bin/expect", "/home/app/entrypoint.sh"]
