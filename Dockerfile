# This Dockerfile attempts to install dependencies, run the tests and build the pocket-core binary
# The result of this Dockerfile will put the pocket-core executable in the $GOBIN/bin, which in turn
# is part of the $PATH

# Base image
FROM golang:1.12-alpine

# Install dependencies
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

# Environment and system dependencies setup
ENV POCKET_PATH=/go/src/github.com/pokt-network/pocket-core/
ENV GO111MODULE="on"

# Create node root directory
RUN mkdir -p ${POCKET_PATH}
WORKDIR $POCKET_PATH

#Install node dependencies
COPY go.mod go.sum ./
RUN go mod download

# Install rest of source code
COPY . .

# Run tests
RUN go test ./tests/...

# Install project dependencies and builds the binary
RUN go build -o ${GOBIN}/bin/pocket-core ./cmd/pocket_core/main.go

# Create app user and add permissions
RUN addgroup -S app \
	&& adduser -S -G app app
RUN chown -R app /go

# Setup the WORKDIR with app user
USER app

# Entrypoint
ENTRYPOINT ["pocket-core"]
