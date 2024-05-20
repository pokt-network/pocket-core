# Based on a previous implementation to make sure we don't break existing deployments.
# https://github.com/pokt-network/pocket-core-deployments/blob/staging/docker/Dockerfile

FROM golang:1.21-alpine as build
RUN apk add --no-cache ca-certificates
WORKDIR /build
ADD . .
RUN go build -o pocket app/cmd/pocket_core/main.go

FROM alpine
RUN apk add --update --no-cache expect bash leveldb-dev tzdata curl \
    && cp /usr/share/zoneinfo/America/New_York /etc/localtime \
    && addgroup --gid 1001 -S app \
    && adduser --uid 1005 -S -G app app

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=build /build/pocket /bin/pocket
COPY .github/workflows/entrypoint.sh /home/app/entrypoint.sh
# allow users to run the container as root and fix the permissions in an easy way
# this could be used as a one time run entrypoint.
# e.g:
# docker-compose:
#   entrypoint: [ "/home/app/change_datadir_ownership_to_app.sh", "/home/app/.pocket"]
#   user: root
# docker:
#   --entrypoint="/home/app/change_datadir_ownership_to_app.sh /home/app/.pocket" -u root
COPY .github/workflows/change_datadir_ownership_to_app.sh /home/app/change_datadir_ownership_to_app.sh

RUN chmod +x /home/app/entrypoint.sh && \
    chmod +x /home/app/change_datadir_ownership_to_app.sh && \
    chown -R app /bin/pocket \
    && mkdir -p /home/app/.pocket/config \
    && chown -R app /home/app/.pocket

# run the container as app user instead of root
USER app

ENTRYPOINT ["/usr/bin/expect", "/home/app/entrypoint.sh"]
