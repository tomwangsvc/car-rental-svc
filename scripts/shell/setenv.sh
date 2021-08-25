#!/bin/bash

# This script:
# -> Sets the environment for working with a GCP project

# Usage:
# -> . setenv.sh

# This file is 'sourced' into the current shell process
# -> Don't set 'nounset', 'errexit', 'pipefail' because that is the responsibility of whatever started the shell process

SCRIPT_NAME="setenv.sh"

echo "${SCRIPT_NAME} -> START at `date '+%Y-%m-%d %H:%M:%S'`..."

log_env_variables ()
{
  echo
  echo "ENV=${1}"
  echo "GCP_PROJECT_ID=${2}"
  echo "GCP_PROJECT_NUMBER=${3}"
  echo
}

if [ $# -lt 1 ]; then
  echo
  echo "${SCRIPT_NAME} -> ERROR: Usage: ${SCRIPT_NAME} dev|prd|stg|uat"
  echo

  else
    export ENV=${1}

    case ${ENV} in
      "dev")
        export GCP_PROJECT_ID="data-fabric-323905"
        export GCP_PROJECT_NUMBER="518937487179"
        log_env_variables ${ENV} ${GCP_PROJECT_ID} ${GCP_PROJECT_NUMBER}
      ;;
      "prd")
        export GCP_PROJECT_ID=""
        export GCP_PROJECT_NUMBER=""
        log_env_variables ${ENV} ${GCP_PROJECT_ID} ${GCP_PROJECT_NUMBER}
      ;;
      "stg")
        export GCP_PROJECT_ID=""
        export GCP_PROJECT_NUMBER=""
        log_env_variables ${ENV} ${GCP_PROJECT_ID} ${GCP_PROJECT_NUMBER}
      ;;
      "uat")
        export GCP_PROJECT_ID=""
        export GCP_PROJECT_NUMBER=""
        log_env_variables ${ENV} ${GCP_PROJECT_ID} ${GCP_PROJECT_NUMBER}
      ;;
      *)
        echo
        echo "${SCRIPT_NAME} -> ERROR: Unsupported environment -> ${ENV}"
        echo
    esac
fi

echo "${SCRIPT_NAME} -> END `date '+%Y-%m-%d %H:%M:%S'`"
