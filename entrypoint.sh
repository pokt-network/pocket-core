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
# POCKET_CORE_SERVICE_WHITELIST = Array list of service whitelist configuration ex. ('["SERVICE1"]'), defaults to []
# POCKET_CORE_DEVELOPER_WHITELIST = Array list of developer whitelist configuration ex. (["DEVELOPER1"]'), defaults to []
# POCKET_CORE_DISPATCH_IP = Dispatch node address (defaults to 127.0.0.1)

cmd="$@"

# Loading pocket-core configurations
## If variable POCKET_CORE_S3_CONFIG_URL is not provided, we use POCKET_CORE_SERVICE_WHITELIST and POCKET_CORE_DEVLEOPER_WHITELIST variables
## If variable POCKET_CORE_S3_CONFIG_URL is provided, we use only the configuration inside the S3 path

if [  ${POCKET_CORE_S3_CONFIG_URL:-false} == false  ]; then
 	echo 'POCKET_CORE_S3_CONFIG_URL env variable not found, Using default configurations'
	if [ ! -f ${POCKET_PATH_DATADIR:-datadir}/service_whitelist.json ]; then
		echo ${POCKET_CORE_SERVICE_WHITELIST:-[]} >  ${POCKET_PATH_DATADIR:-datadir}/service_whitelist.json
	fi

	if [ ! -f ${POCKET_PATH_DATADIR:-datadir}/developer_whitelist.json ]; then
		echo ${POCKET_CORE_DEVELOPER_WHITELIST:-[]} > ${POCKET_PATH_DATADIR:-datadir}/developer_whitelist.json
	fi

	if [ ! -f ${POCKET_PATH_DATADIR:-datadir}/chains.json ]; then
		echo ${POCKET_CORE_CHAINS:-[]} > ${POCKET_PATH_DATADIR:-datadir}/chains.json
	fi
else
 	echo 'Downloading node configurations from S3'
	aws s3 sync $POCKET_CORE_S3_CONFIG_URL ${POCKET_PATH_DATADIR:-datadir}
fi

# Running tests

if [ ${POCKET_CORE_UNIT_TESTS:-false} == true ]; then
    echo "Initializing unit testing"
    go test ./tests/unit/...
fi


exec $cmd
