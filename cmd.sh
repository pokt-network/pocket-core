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
# POCKET_CORE_DISPATCH_IP = Dispatch node address (defaults to 127.0.0.1)
# POCKET_CORE_DISPATCH_PORT = Dispatch node port (defaults to 8081)

# Dispatch only
# AWS_DYNAMODB_ENDPOINT = AWS dynamodb endpoint (defaults to dynamodb.us-east-1.amazonaws.com)
# AWS_DYNAMODB_TABLE = AWS dynamodb table name (defaults to dispatch-peers-staging)
# AWS_DYNAMODB_REGION = AWS dynamodb region (defaults to us-east-1)

# Service only
# POCKET_CORE_SERVICE_GID = GID of the service node (required)
# POCKET_CORE_SERVICE_IP = IP of the Pocket Core service node (required)
# POCKET_CORE_SERVICE_PORT = Port of the Pocket Core service node (required)

# Start pocket-core
if [ $POCKET_CORE_NODE_TYPE = "dispatch" ]; then
	echo 'Starting pocket-core dispatch'

	exec pocket-core --dispatch \
		--datadirectory ${POCKET_PATH_DATADIR:-datadir} \
		--dbend ${POCKET_CORE_AWS_DYNAMODB_ENDPOINT:-dynamodb.us-east-1.amazonaws.com} \
		--dbtable ${POCKET_CORE_AWS_DYNAMODB_TABLE:-dispatch-peers-staging} \
		--dbregion ${POCKET_CORE_AWS_DYNAMODB_REGION:-us-east-1} \
		--disip ${POCKET_CORE_DISPATCH_IP:-127.0.0.1} \
		--disrport ${POCKET_CORE_DISPATCH_PORT:-8081}

elif [ $POCKET_CORE_NODE_TYPE = "service" ]; then
	echo 'Starting pocket-core service'

	exec pocket-core --datadirectory ${POCKET_PATH_DATADIR:-datadir} \
		--disip ${POCKET_CORE_DISPATCH_IP:-127.0.0.1} \
		--gid ${POCKET_CORE_SERVICE_GID:-GID2} \
		--ip ${POCKET_CORE_SERVICE_IP:-127.0.0.1} \
		--disrport ${POCKET_CORE_DISPATCH_PORT:-8081} \
		--port ${POCKET_CORE_SERVICE_PORT:-8081}

else
	echo 'Need to specify a node type, either dispatch or service.'
	exit 1
fi
