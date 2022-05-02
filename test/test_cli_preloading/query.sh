#!/usr/bin/env bash
set -Eeuo pipefail

USERNAME=${1:-test}
PASSWORD=${2:-test}

test1 () {
  docker-compose -f docker-compose.tools.yml run medco-cli-client --user $USERNAME --password $PASSWORD --o /data/result.csv $1 $2
  result="$(awk -F "\"*,\"*" '{print $2}' ../result.csv)"
  if [ "${result}" != "${3}" ];
  then
  echo "$1 $2: test failed"
  echo "result: ${result}" && echo "expected result: ${3}"
  exit 1
  fi
}

test2 () {
  docker-compose -f docker-compose.tools.yml run medco-cli-client --user $USERNAME --password $PASSWORD --o /data/result.csv $1 $2
  result="$(awk -F "\"*,\"*" '{print $2}' ../result.csv)"
  if [ "${result}" == "${3}" ];
  then
  echo "$1 $2: WARNING - result is the same"
  echo "result: ${result}" && echo "expected result: ${3}"
  fi
}

test3 () {
  docker-compose -f docker-compose.tools.yml run medco-cli-client --user $USERNAME --password $PASSWORD --o /data/result.csv query enc::1
  queryIDs1="$(awk -F "\"*,\"*" '{if (NR != 1) {print}}' ../result.csv | sort | awk -F "\"*,\"*" '{print $4}' | awk 'BEGIN{ORS=","}1' | sed 's/.$//')"
  echo "queryIDs1 $queryIDs1"
  docker-compose -f docker-compose.tools.yml run medco-cli-client --user $USERNAME --password $PASSWORD add-saved-cohorts -c testCohortQuery1 -q $queryIDs1

  docker-compose -f docker-compose.tools.yml run medco-cli-client --user $USERNAME --password $PASSWORD --o /data/result.csv query enc::2
  queryIDs2="$(awk -F "\"*,\"*" '{if (NR != 1) {print}}' ../result.csv | sort | awk -F "\"*,\"*" '{print $4}' | awk 'BEGIN{ORS=","}1' | sed 's/.$//')"
  echo "queryIDs2 $queryIDs2"
  docker-compose -f docker-compose.tools.yml run medco-cli-client --user $USERNAME --password $PASSWORD add-saved-cohorts -c testCohortQuery2 -q $queryIDs2

  docker-compose -f docker-compose.tools.yml run medco-cli-client --user $USERNAME --password $PASSWORD --o /data/result.csv query enc::1 OR enc::2
  result1="$(awk -F "\"*,\"*" '{print $3}' ../result.csv)"
  docker-compose -f docker-compose.tools.yml run medco-cli-client --user $USERNAME --password $PASSWORD --o /data/result.csv query enc::1 AND enc::2
  result2="$(awk -F "\"*,\"*" '{print $3}' ../result.csv)"
  docker-compose -f docker-compose.tools.yml run medco-cli-client --user $USERNAME --password $PASSWORD --o /data/result.csv query chr::testCohortQuery1 OR enc::2
  resultWithPsID1="$(awk -F "\"*,\"*" '{print $3}' ../result.csv)"
  if [ "${result1}" != "${resultWithPsID1}" ];
  then
  echo "test 3 failed"
  echo "result: ${resultWithPsID1}" && echo "expected result: ${result1}"
  exit 1
  fi
  docker-compose -f docker-compose.tools.yml run medco-cli-client --user $USERNAME --password $PASSWORD --o /data/result.csv query enc::1 AND chr::testCohortQuery2
  resultWithPsID2="$(awk -F "\"*,\"*" '{print $3}' ../result.csv)"
  if [ "${result2}" != "${resultWithPsID2}" ];
  then
  echo "test 3 failed"
  echo "result: ${resultWithPsID2}" && echo "expected result: ${result2}"
  exit 1
  fi
}

function timing() { echo "query clr::/E2ETEST/SPHNv2020.1/DeathStatus/ OR clr::/E2ETEST/SPHNv2020.1/DeathStatus/ ${1} AND clr::/E2ETEST/SPHNv2020.1/DeathStatus/:/E2ETEST/DeathStatus-status/death/:/SPHNv2020.1/DeathStatus/ ${2} AND clr::/E2ETEST/I2B2/Demographics/Gender/Female/ OR clr::/E2ETEST/I2B2/Demographics/Gender/Male/ ${3} -t ${4}"; };
timingResultNonZeroExpected="$(printf -- "count\n165\n165\n165")"
timingResultZeroExpected="$(printf -- "count\n0\n0\n0")"
test4() {
  docker-compose -f docker-compose.tools.yml run medco-cli-client --user $USERNAME --password $PASSWORD --o /data/result.csv \
  $(timing $1 $2 $3 $4)
  result="$(awk -F "\"*,\"*" '{print $2}' ../result.csv)"
  if [ "${result}" != "${5}" ];
  then
  echo "query timing ${1} ${2} ${3} ${4}: test failed"
  echo "result: ${result}" && echo "expected result: ${5}"
  exit 1
  fi

}


test5() {
  docker-compose -f docker-compose.tools.yml run medco-cli-client --user $USERNAME --password $PASSWORD --o /data/result.csv \
  query -s "${2}" "${1}"
  result="$(awk -F "\"*,\"*" '{print $2}' ../result.csv)"
  if [ "${result}" != "${3}" ];
  then
  echo "query temporal sequence ${2}: test failed"
  echo "result: ${result}" && echo "expected result: ${3}"
  exit 1
  fi
}

query1="enc::1 OR enc::2 AND enc::3"
resultQuery1="$(printf -- "count\n1\n1\n1")"

query2="clr::/E2ETEST/e2etest/1/ OR enc::2 AND enc::3"
resultQuery2="$(printf -- "count\n1\n1\n1")"

query3="clr::/E2ETEST/e2etest/1/"
resultQuery3="$(printf -- "count\n3\n3\n3")"

query4="clr::/E2ETEST/e2etest/1/:/E2ETEST/modifiers/1/:/e2etest/%"
resultQuery4="$(printf -- "count\n2\n2\n2")"

echo "enc::1
enc::2" > deployments/query_file.txt
query5="enc::3 AND file::/data/query_file.txt"
resultQuery5a="$(printf -- "count\n1\n1\n1")"
resultQuery5b="$(printf -- "count\n3\n3\n3")"

query6="clr::/E2ETEST/e2etest/1/::EQ:NUMBER:10"
resultQuery6="$(printf -- "count\n1\n1\n1")"

query7="clr::/E2ETEST/e2etest/1/:/E2ETEST/modifiers/1/:/e2etest/1/::EQ:NUMBER:15"
resultQuery7="$(printf -- "count\n1\n1\n1")"

query8="clr::/E2ETEST/e2etest/1/::BETWEEN:NUMBER:5 and 25"
resultQuery8="$(printf -- "count\n2\n2\n2")"

query9="enc::1 OR clr::/E2ETEST/e2etest/2/::GE:NUMBER:25 AND clr::/E2ETEST/e2etest/2/:/E2ETEST/modifiers/2/:/e2etest/2/::LT:NUMBER:21"
resultQuery9="$(printf -- "count\n2\n2\n2")"

query10="clr::/E2ETEST/e2etest/2/:/E2ETEST/modifiers/2text/:/e2etest/2/::IN:TEXT:'abc','de'"
resultQuery10="$(printf -- "count\n2\n2\n2")"

query11="clr::/E2ETEST/e2etest/3/:/E2ETEST/modifiers/3text/:/e2etest/3/::LIKE[begin]:TEXT:ab"
resultQuery11="$(printf -- "count\n2\n2\n2")"

query12="clr::/E2ETEST/e2etest/2/:/E2ETEST/modifiers/2text/:/e2etest/2/::LIKE[contains]:TEXT:cd"
resultQuery12="$(printf -- "count\n1\n1\n1")"

query13="clr::/E2ETEST/e2etest/3/:/E2ETEST/modifiers/3text/:/e2etest/3/::LIKE[end]:TEXT:bc"
resultQuery13="$(printf -- "count\n0\n0\n0")"

query14="clr::/SPHN/SPHNv2020.1/FophDiagnosis/ OR clr::/SPHN/SPHNv2020.1/DeathStatus/ WITH clr::/SPHN/SPHNv2020.1/FophDiagnosis/ THEN clr::/SPHN/SPHNv2020.1/DeathStatus/ THEN clr::/SPHN/SPHNv2020.1/DeathStatus/"
resultQuery14a="$(printf -- "count\n228\n228\n228")"

resultQuery14b="$(printf -- "count\n0\n0\n0")"

pushd deployments/dev-local-3nodes/
echo "Testing query with test user..."

test1 "query " "${query1}" "${resultQuery1}"
test1 "query " "${query2}" "${resultQuery2}"
test1 "query " "${query3}" "${resultQuery3}"
test1 "query " "${query4}" "${resultQuery4}"
test1 "query " "${query5}" "${resultQuery5a}"
test1 "query " "${query6}" "${resultQuery6}"
test1 "query " "${query7}" "${resultQuery7}"
test1 "query " "${query8}" "${resultQuery8}"
test1 "query " "${query9}" "${resultQuery9}"
test1 "query " "${query10}" "${resultQuery10}"
test1 "query " "${query11}" "${resultQuery11}"
test1 "query " "${query12}" "${resultQuery12}"
test1 "query " "${query13}" "${resultQuery13}"

echo "Testing query with test_explore_patient_list user..."
USERNAME="${1:-test}_explore_patient_list"

test1 "query " "${query1}" "${resultQuery1}"
test1 "query " "${query2}" "${resultQuery2}"
test1 "query " "${query3}" "${resultQuery3}"
test1 "query " "${query4}" "${resultQuery4}"
test1 "query " "${query5}" "${resultQuery5a}"

USERNAME="${1:-test}_explore_count_global"
test1 "query " "${query5}" "${resultQuery5b}"

USERNAME="${1:-test}_explore_count_global_obfuscated"
test2 "query " "${query5}" "${resultQuery5b}"

echo "Testing query with cohort as query term..."
USERNAME="${1:-test}"
test3

echo "Testing query with timing settings features..."

test4 "any" "any" "any" "any" "${timingResultNonZeroExpected}"
test4 "sameinstancenum" "sameinstancenum" "sameinstancenum" "sameinstancenum" "${timingResultZeroExpected}"
test4 "samevisit" "samevisit" "samevisit" "samevisit" "${timingResultZeroExpected}"
test4 "sameinstancenum" "sameinstancenum" "any" "sameinstancenum" "${timingResultNonZeroExpected}"
test4 "samevisit" "samevisit" "any" "samevisit" "${timingResultNonZeroExpected}"

echo "Testing query with event sequences features..."

test5 "${query14}" "before,first,startdate,first,startdate:sametime,first,startdate,first,startdate" "${resultQuery14a}"
test5 "${query14}" "sametime,first,startdate,first,startdate,moreorequal,23,days,more,12,days,less,30,hours:sametime,first,startdate,first,startdate,moreorequal,20,days,more,11,days,more,23,hours" "${resultQuery14b}"

popd
exit 0