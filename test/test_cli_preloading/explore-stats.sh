#!/usr/bin/env bash
set -Eeuo pipefail

USERNAME=${1:-test}
PASSWORD=${2:-test}


testBioref() {
  echo "launching test for bioref features"
  nbBuckets=$1
  expected="$2"

  logFile="/data/result.csv"

  docker-compose -f docker-compose.tools.yml run medco-cli-client \
    --user $USERNAME --password $PASSWORD  --o "$logFile"  \
    explore-stats clr::/E2ETEST/e2etest/bioref/ testCohortBioref -b "${nbBuckets}"

  result="$(cat ../result.csv)"
  if [ "${result}" != "${expected}" ]; then
    echo "result: ${result}" && echo "expected result: ${expected}"
    exit 1
  fi

  echo successful bioref test
}


testBioref 3 "${expectedBioref1}"