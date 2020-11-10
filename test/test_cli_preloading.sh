#!/usr/bin/env bash
set -Eeuo pipefail

USERNAME=${1:-test}
PASSWORD=${2:-test}


# test1
searchConceptChildren1="/"
resultSearchConceptChildren1="PATH  TYPE
                              /E2ETEST/e2etest/concept_container
                              /I2B2/I2B2/ concept_container
                              /SPHN/SPHNv2020.1/  concept_container"

searchConceptChildren2="/E2ETEST/e2etest/"
resultSearchConceptChildren2="PATH  TYPE
                              /E2ETEST/e2etest/1/ concept
                              /E2ETEST/e2etest/2/ concept
                              /E2ETEST/e2etest/3/ concept
                              /E2ETEST/modifiers/ modifier_folder"

searchModifierChildren="/E2ETEST/modifiers/ /e2etest/% /E2ETEST/e2etest/1/"
resultSearchModifierChildren="PATH  TYPE
                              /E2ETEST/modifiers/1/ modifier"

# test 2
query1="enc::1 OR enc::2 AND enc::3"
resultQuery1="$(printf -- "count\n1\n1\n1")"

query2="clr::/E2ETEST/e2etest/1/ OR enc::2 AND enc::3"
resultQuery2="$(printf -- "count\n1\n1\n1")"

query3="clr::/E2ETEST/e2etest/1/"
resultQuery3="$(printf -- "count\n2\n2\n2")"

query4="clr::/E2ETEST/e2etest/1/::/E2ETEST/modifiers/1/:/e2etest/%"
resultQuery4="$(printf -- "count\n1\n1\n1")"

echo "enc::1
enc::2" > deployments/query_file.txt
query5="enc::3 AND file::/data/query_file.txt"
resultQuery5="$(printf -- "count\n1\n1\n1")"

# test3
getSavedCohortHeaders="node_index,cohort_name,cohort_id,query_id,creation_date,update_date"
getSavedCohort1="$(printf -- "node_index cohort_name cohort_id query_id\n0 testCohort -1 -1\n1 testCohort -1 -1\n2 testCohort -1 -1")"
getSavedCohort2="$(printf -- "node_index cohort_name query_id\n0 testCohort -1\n0 testCohort2 -1\n1 testCohort -1\n1 testCohort2 -1\n2 testCohort -1\n2 testCohort2 -1")"

# test4
timerHeaders="node_index,timer_description,duration_milliseconds"
survivalDays="$(printf -- "time_granularity,node_index,group_id,initial_count,time_point,event_of_interest_count,censoring_event_count\nday,0,Full cohort,684,0,0,0\nday,0,Full cohort,684,1,0,0\nday,0,Full cohort,684,2,0,0\nday,0,Full cohort,684,3,0,0\nday,0,Full cohort,684,4,0,0\nday,0,Full cohort,684,5,3,0")"
survivalWeeks="$(printf -- "time_granularity,node_index,group_id,initial_count,time_point,event_of_interest_count,censoring_event_count\nweek,0,Full cohort,684,0,0,0\nweek,0,Full cohort,684,1,3,0\nweek,0,Full cohort,684,2,18,0\nweek,0,Full cohort,684,3,3,0\nweek,0,Full cohort,684,4,3,0\nweek,0,Full cohort,684,5,6,0")"
survivalMonths="$(printf -- "time_granularity,node_index,group_id,initial_count,time_point,event_of_interest_count,censoring_event_count\nmonth,0,Full cohort,684,0,0,0\nmonth,0,Full cohort,684,1,30,0\nmonth,0,Full cohort,684,2,21,0\nmonth,0,Full cohort,684,3,30,0\nmonth,0,Full cohort,684,4,30,6\nmonth,0,Full cohort,684,5,30,0")"
survivalYears="$(printf -- "time_granularity,node_index,group_id,initial_count,time_point,event_of_interest_count,censoring_event_count\nyear,0,Full cohort,684,0,0,0\nyear,0,Full cohort,684,1,363,126\nyear,0,Full cohort,684,2,114,42\nyear,0,Full cohort,684,3,18,21\nyear,0,Full cohort,684,4,0,0\nyear,0,Full cohort,684,5,0,0")"

# test5
survivalSubGroup1="$(printf -- "time_granularity,node_index,group_id,initial_count,time_point,event_of_interest_count,censoring_event_count\nweek,0,Female,414,0,0,0\nweek,0,Female,414,1,0,0\nweek,0,Female,414,2,18,0\nweek,0,Female,414,3,3,0\nweek,0,Female,414,4,3,0\nweek,0,Female,414,5,6,0")"
survivalSubGroup2="$(printf -- "week,0,Male,270,0,0,0\nweek,0,Male,270,1,3,0\nweek,0,Male,270,2,0,0\nweek,0,Male,270,3,0,0\nweek,0,Male,270,4,0,0\nweek,0,Male,270,5,0,0")"

test1 () {
  docker-compose -f docker-compose.tools.yml run medco-cli-client --user $USERNAME --password $PASSWORD --o /data/result.csv $1 $2
  result="$(cat ../result.csv | tr -d '\r\n\t ')"
  expectedResult="$(echo "${3}" | tr -d '\r\n\t ')"
  if [ "${result}" != "${expectedResult}" ];
  then
  echo "$1 $2: test failed"
  echo "result: ${result}" && echo "expected result: ${expectedResult}"
  exit 1
  fi
}

test2 () {
  docker-compose -f docker-compose.tools.yml run medco-cli-client --user $USERNAME --password $PASSWORD --o /data/result.csv $1 $2
  result="$(awk -F "\"*,\"*" '{print $2}' ../result.csv)"
  if [ "${result}" != "${3}" ];
  then
  echo "$1 $2: test failed"
  echo "result: ${result}" && echo "expected result: ${3}"
  exit 1
  fi
}

test3 () {
  docker-compose -f docker-compose.tools.yml run medco-cli-client --user $USERNAME --password $PASSWORD --o /data/result.csv get-saved-cohorts
  result="$(awk -F',' 'NR==1{print $0}' ../result.csv)"
  if [ "${result}" != "${getSavedCohortHeaders}" ];
  then
  echo "get-saved-cohorts headers: test failed"
  echo "result: ${result}" && echo "expected result: ${getSavedCohortHeaders}"
  exit 1
  fi

  result="$(awk -F',' '{print $1,$2,$3,$4}' ../result.csv)"
  if [ "${result}" != "${getSavedCohort1}" ];
  then
  echo "get-saved-cohorts content before update: test failed"
  echo "result: ${result}" && echo "expected result: ${getSavedCohort1}"
  exit 1
  fi

  docker-compose -f docker-compose.tools.yml run medco-cli-client --user $USERNAME --password $PASSWORD add-saved-cohorts -c testCohort2 -p $(echo -1,-1,-1)
  docker-compose -f docker-compose.tools.yml run medco-cli-client --user $USERNAME --password $PASSWORD --o /data/result.csv get-saved-cohorts
  result="$(awk -F',' '{print $1,$2,$4}' ../result.csv)"
  if [ "${result}" != "${getSavedCohort2}" ];
  then
  echo "get-saved-cohorts content after added new cohort: test failed"
  echo "result: ${result}" && echo "expected result: ${getSavedCohort2}"
  exit 1
  fi


  docker-compose -f docker-compose.tools.yml run medco-cli-client --user $USERNAME --password $PASSWORD update-saved-cohorts -c testCohort2 -p $(echo -1,-1,-1)
  docker-compose -f docker-compose.tools.yml run medco-cli-client --user $USERNAME --password $PASSWORD --o /data/result.csv get-saved-cohorts
  result="$(awk -F',' '{print $1,$2,$4}' ../result.csv)"
  if [ "${result}" != "${getSavedCohort2}" ];
  then
  echo "get-saved-cohorts content after update: test failed"
  echo "result: ${result}" && echo "expected result: ${getSavedCohort2}"
  exit 1
  fi

  docker-compose -f docker-compose.tools.yml run medco-cli-client --user $USERNAME --password $PASSWORD remove-saved-cohorts -c testCohort2
  docker-compose -f docker-compose.tools.yml run medco-cli-client --user $USERNAME --password $PASSWORD --o /data/result.csv get-saved-cohorts
  result="$(awk -F',' '{print $1,$2,$3,$4}' ../result.csv)"
  if [ "${result}" != "${getSavedCohort1}" ];
  then
  echo "get-saved-cohorts content after update: test failed"
  echo "result: ${result}" && echo "expected result: ${getSavedCohort1}"
  exit 1
  fi

}

test4 () {
  docker-compose -f docker-compose.tools.yml run medco-cli-client --user $USERNAME --password $PASSWORD -o /data/result.csv srva  -c testCohort -l 6 -g ${1}  -s /SPHN/SPHNv2020.1/FophDiagnosis/ -e /SPHN/SPHNv2020.1/DeathStatus/ -y 126:1 -d /data/timers.csv
  
  result="$(awk -F',' 'NR==1{print $0}' ../timers.csv)"
  if [ "${result}" != "${timerHeaders}" ];
  then
  echo "timer headers $1: test failed"
  echo "result: ${result}" && echo "expected result: ${timerHeaders}"
  exit 1
  fi
  
  result="$(awk -F',' 'NR==1, NR==7 {print $0}' ../result.csv)"
  if [ "${result}" != "${2}" ];
  then
  echo "survival analysis $1: test failed"
  echo "result: ${result}" && echo "expected result: ${2}"
  exit 1
  fi

}

test5 () {
  docker-compose -f docker-compose.tools.yml run \
    -v "${PWD}/../../test/survival_test_parameters.yaml":/parameters/survival_test_parameters.yaml \
    medco-cli-client --user $USERNAME --password $PASSWORD -o /data/result.csv srva -d /data/timers.csv \
    -p /parameters/survival_test_parameters.yaml

  result="$(awk -F',' 'NR==1, NR==7 {print $0}' ../result.csv)"
  if [ "${result}" != "${1}" ];
  then
  echo "survival analysis sub group 1: test failed"
  echo "result: ${result}" && echo "expected result: ${1}"
  exit 1
  fi

  result="$(awk -F',' 'NR==8, NR==13 {print $0}' ../result.csv)"
  if [ "${result}" != "${2}" ];
  then
  echo "survival analysis sub group 2: test failed"
  echo "result: ${result}" && echo "expected result: ${2}"
  exit 1
  fi
}

pushd deployments/dev-local-3nodes/
echo "Testing concept-children..."

test1 "concept-children" "${searchConceptChildren1}" "${resultSearchConceptChildren1}"
test1 "concept-children" "${searchConceptChildren2}" "${resultSearchConceptChildren2}"

echo "Testing modifier-children..."

test1 "modifier-children" "${searchModifierChildren}" "${resultSearchModifierChildren}"

echo "Testing query..."

test2 "query patient_list" "${query1}" "${resultQuery1}"
test2 "query patient_list" "${query2}" "${resultQuery2}"
test2 "query patient_list" "${query3}" "${resultQuery3}"
test2 "query patient_list" "${query4}" "${resultQuery4}"
test2 "query patient_list" "${query5}" "${resultQuery5}"

echo "Testing saved-cohorts features..."

test3

test4 "day" "${survivalDays}"
test4 "week" "${survivalWeeks}"
test4 "month" "${survivalMonths}"
test4 "year" "${survivalYears}"

test5 "${survivalSubGroup1}" "${survivalSubGroup2}"

echo "CLI test 1/2 successful!"
popd
exit 0