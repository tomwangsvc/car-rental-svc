#!/bin/bash

# This script:
# -> Builds, tests and starts an service

# Exit script if you try to use an uninitialized variable.
set -o nounset

# Exit script if a statement returns a non-true return value.
set -o errexit

# Use the error status of the first failure, rather than that of the last item in a pipeline.
set -o pipefail

SCRIPT_DIR="$( cd "$(dirname "${0}")" ; pwd -P )"
SCRIPT_DIR_NAME=${SCRIPT_DIR##*/}
SCRIPT_NAME=`basename ${0}`
#echo "SCRIPT_DIR=${SCRIPT_DIR}"
#echo "SCRIPT_DIR_NAME=${SCRIPT_DIR_NAME}"
#echo "SCRIPT_NAME=${SCRIPT_NAME}"

SVC=${SCRIPT_DIR_NAME}

echo "${SCRIPT_NAME} -> START at `date '+%Y-%m-%d %H:%M:%S'`..."

if [ $# -lt 1 ]; then
    export ENV="dev"
  else
    export ENV=${1}
fi

case ${ENV} in
  "dev")
    export GCP_PROJECT_ID="data-fabric-323905"
    export GCP_PROJECT_NUMBER="518937487179"
  ;;
  *)
    echo "${SCRIPT_NAME} -> ERROR: Unsupported environment -> ${ENV}"
    echo "${SCRIPT_NAME} -> END `date '+%Y-%m-%d %H:%M:%S'`"
    exit 1
esac

echo
echo "SVC=${SVC}"
echo "ENV=${ENV}"
echo "GCP_PROJECT_ID=${GCP_PROJECT_ID}"
echo "GCP_PROJECT_NUMBER=${GCP_PROJECT_NUMBER}"
echo

cd ${SCRIPT_DIR}

echo "Formatting..."
set +e
go fmt ${SVC}/...
set -e
echo

echo "Testing..."
go test -timeout 5000ms -cover -mod=vendor ${SVC}/...
echo

echo "Cleaning bin..."
[ -d ./bin ] && rm -rf ./bin
mkdir bin
[ -d ./schema ] && cp -r ./schema/ ./bin
[ -d ./templates ] && cp -r ./templates/ ./bin
echo

echo "Building..."
go build -mod=vendor -o ./bin/${SVC} ./cmd/service
echo

echo "Running..."
./bin/${SVC}
