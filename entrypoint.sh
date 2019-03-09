#!bin/bash
# Set utils
set -o errexit
set -o pipefail
set -o nounset

# Environment variables
# All nodes
# POCKET_CORE_NODE_TYPE = dispatch | service
# POCKET_PATH_DATADIR = absolute path to the datadirectory
# POCKET_CORE_S3_CONFIG_URL = S3 directory to download configurations from
# AWS_ACCESS_KEY_ID = aws access key to download the S3 / connect to dynamodb (dispatch only)
# AWS_SECRET_ACCESS_KEY = aws secret key to download the S3 / connect to dynamodb (dispatch only)

# Dispatch only
# AWS_DYNAMODB_ENDPOINT = AWS dynamodb endpoint (defaults to dynamodb.us-east-1.amazonaws.com)
# AWS_DYNAMODB_TABLE = AWS dynamodb table name (defaults to dispatch-peers-staging)
# AWS_DYNAMODB_REGION = AWS dynamodb region (defaults to us-east-1)

# Service only
# POCKET_CORE_DISPATCH_IP = Dispatch node address (defaults to 127.0.0.1)

cmd="$@"

# Download node configurations from S3

if [  ${POCKET_CORE_S3_CONFIG_URL:-false} != false ]; then
    echo 'Downloading node configurations'
    aws s3 sync $POCKET_CORE_S3_CONFIG_URL ${POCKET_PATH_DATADIR:-datadir}
fi

# Running tests

if [ ${POCKET_CORE_UNIT_TESTS:-false} == true ]; then
    echo "Initializing unit testing"
    go test ./tests/unit/...  
fi

if [ ${POCKET_CORE_INTEGRATION_TESTS:-false}  == true ]; then
    echo "Initializing integration testing"
    go test ./tests/integration/... --disip ${POCKET_CORE_DISPATCH_IP:-127.0.0.1} --disrport ${POCKET_CORE_DISPATCH_PORT:-8081}
fi


exec $cmd
