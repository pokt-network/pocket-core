FROM golang:1.11-alpine

ENV POCKET_PATH=/go/src/github.com/pokt-network/pocket-core/
RUN mkdir -p ${POCKET_PATH}
COPY . $POCKET_PATH

RUN apk add curl git && \
	rm -rf /var/cache/apk/* && \
	curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh


RUN cd ${POCKET_PATH} &&  \
	dep ensure && \
	go build ./cmd/pocket_core/main.go && \
	ls $GOBIN && \
	ls /go/src/github.com/ && \
	ls $GOROOT

RUN mkdir $POCKET_PATH/.pocket && echo "[]" >> $POCKET_PATH/.pocket/chains.json

# TODO: Run tests
# TODO:


RUN addgroup -S app \
	&& adduser -S -G app app

RUN chown -R app /go

USER app
WORKDIR $POCKET_PATH

CMD "go run cmd/pocket_core/main.go --dispatch"
