#!/bin/sh
datadir=$1
echo "Attempting to fix ${datadir} permissions to be own by app:app"
chown -R app:app $datadir
echo "${datadir} permissions applied."
echo "Please turn off entrypoint override and ensure you are using user app or user 1005 when start container."
exit 0