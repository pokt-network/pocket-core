# This Dockerfile attempts to install dependencies, run the tests and build the pocket-core binary
# The result of this Dockerfile will put the pocket-core executable in the $GOBIN/bin, which in turn
# is part of the $PATH

# Base image
FROM golang:1.11-alpine as builder

# Environment and system dependencies setup
ENV POCKET_PATH=/go/src/github.com/pokt-network/pocket-core/
RUN mkdir -p ${POCKET_PATH}
COPY . $POCKET_PATH

RUN apk -v --update --no-cache add \
		git \
		groff \
		less \
		mailcap \
		dep \
                gcc \
                libc-dev \
		bash && \
	rm /var/cache/apk/* || true

# Setup data directory
RUN mkdir ${POCKET_PATH}datadir
ENV POCKET_PATH_DATADIR=${POCKET_PATH}datadir

# Install project dependencies and builds the binary
RUN cd ${POCKET_PATH} && dep ensure && go build -o ${GOBIN}/bin/pocket-core ./cmd/pocket_core/main.go


# container app
FROM golang:1.11-alpine 

ENV POCKET_PATH=/go/src/github.com/pokt-network/pocket-core/
RUN mkdir -p ${POCKET_PATH}/datadir

COPY --from=builder ${GOBIN}/bin/pocket-core ${GOBIN}/bin/pocket-core   
COPY --from=builder ${POCKET_PATH}tests ${POCKET_PATH}   
COPY entrypoint.sh ${POCKET_PATH}   

RUN apk -v --update --no-cache add \
		curl \
		python \
		py-pip \
		bash && \
	pip install --upgrade --no-cache awscli s3cmd && \
	apk -v --purge del py-pip && \
	rm /var/cache/apk/* || true

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
