#!bin/bash
set -o errexit
set -o pipefail
set -o nounset

# Run integration tests against local env

go test ./tests/integration/...  -url ${POCKET_CORE_DISPATCH_IP:-127.0.0.1} -requestTimeout ${POCKET_CORE_INTEGRATION_TEST_TIMEOUT:-400}
