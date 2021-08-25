#!/bin/bash

# This script:
# -> Builds, tests and starts an api

# Exit script if you try to use an uninitialized variable.
set -o nounset

# Exit script if a statement returns a non-true return value.
set -o errexit

# Use the error status of the first failure, rather than that of the last item in a pipeline.
set -o pipefail

SCRIPT_DIR="$( cd "$(dirname "${0}")" ; pwd -P )"
SCRIPT_DIR_NAME=${SCRIPT_DIR##*/}
SCRIPT_NAME=`basename ${0}`
# echo "SCRIPT_DIR=${SCRIPT_DIR}"
# echo "SCRIPT_DIR_NAME=${SCRIPT_DIR_NAME}"
# echo "SCRIPT_NAME=${SCRIPT_NAME}"

cd ${SCRIPT_DIR}/../..
API=`basename $(pwd)`
# echo "API=${API}"

if [ $# -lt 1 ]; then
    export ENV="dev"
  else
    export ENV=${1}
fi

if [ ${ENV} = "prd" ]; then
  gsutil cp gs://l214552987832909-svc-postman-environments-production/${API}/${ENV}/* ${SCRIPT_DIR}/../../test-api
else
  gsutil cp gs://l214552987832909-svc-postman-environments-non-production/${API}/${ENV}/* ${SCRIPT_DIR}/../../test-api;
fi
echo
