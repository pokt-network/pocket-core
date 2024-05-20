#!/bin/sh
datadir=$1
# Changing the permissions is necessary here because previous versions of our Dockerfile
# did not specify `app` as the user, so Docker defaulted to `root`.
# This script facilitates transition for those who have not specified the `app` user at the start of the container.
# It changes the ownership to the proper user and group (`app:app`), as declared in the Dockerfile.
# The specific ownership by user 'app' and group 'app' is required to ensure that the `app` user
# specified in the Dockerfile will have full access to the relevant directory.
echo "Attempting to fix ${datadir} permissions to be owned by app:app"
chown -R app:app $datadir
echo "${datadir} permissions applied."
echo "Please turn off entrypoint override and ensure you are using user `app` or user `1005` when start container."
exit 0