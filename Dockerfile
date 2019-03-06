# This Dockerfile attempts to install dependencies, run the tests and build the pocket-core binary
# The result of this Dockerfile will put the pocket-core executable in the $GOBIN/bin, which in turn
# is part of the $PATH

# Base image
FROM golang:1.11-alpine

# Environment and system dependencies setup
ENV POCKET_PATH=/go/src/github.com/pokt-network/pocket-core/
RUN mkdir -p ${POCKET_PATH}
COPY . $POCKET_PATH
RUN apk -v --update --no-cache add \
		curl \
		git \
		python \
		py-pip \
		groff \
		less \
		mailcap \
		dep \
                gcc \
                libc-dev \
		bash && \
	pip install --upgrade --no-cache awscli s3cmd python-magic && \
	apk -v --purge del py-pip && \
	rm /var/cache/apk/* || true

# Setup data directory
RUN mkdir ${POCKET_PATH}datadir
ENV POCKET_PATH_DATADIR=${POCKET_PATH}datadir

# Install project dependencies and builds the binary
RUN cd ${POCKET_PATH} && dep ensure && go build -o ${GOBIN}/bin/pocket-core ./cmd/pocket_core/main.go

# TODO: Run tests
#RUN go test tests/unit/...

# Create app user and add permissions
RUN addgroup -S app \
	&& adduser -S -G app app
RUN chown -R app /go

# Setup the WORKDIR with app user
USER app
WORKDIR $POCKET_PATH

# Expose port 8081
EXPOSE 8081

# Entrypoint
ENTRYPOINT ["/bin/bash", "entrypoint.sh" ]

CMD ["/bin/bash", "cmd.sh"]
