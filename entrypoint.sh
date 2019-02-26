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

# Download node configurations from S3
echo 'Downloading node configurations'
aws s3 sync $POCKET_CORE_S3_CONFIG_URL $POCKET_PATH_DATADIR

# Start pocket-core
if [ $POCKET_CORE_NODE_TYPE = "dispatch" ]; then
	echo 'Starting pocket-core dispatch'
	exec pocket-core --dispatch --datadirectory $POCKET_PATH_DATADIR --dbend ${AWS_DYNAMODB_ENDPOINT:-dynamodb.us-east-1.amazonaws.com} --dbtable ${AWS_DYNAMODB_TABLE:-dispatch-peers-staging} --dbregion ${AWS_DYNAMODB_REGION:-us-east-1}
elif [ $POCKET_CORE_NODE_TYPE = "service" ]; then
	echo 'Starting pocket-core service'
	exec pocket-core --datadirectory $POCKET_PATH_DATADIR --disip ${POCKET_CORE_DISPATCH_IP:-127.0.0.1}
else
	echo 'Need to specify a node type, either dispatch or service.'
	exit 1
fi
