#!/usr/bin/env bash
set -Eeuo pipefail

USERNAME=${1:-test}
PASSWORD=${2:-test}


# test1
getSavedCohortHeaders="node_index,cohort_name,cohort_id,query_id,creation_date,update_date"
getSavedCohort1="$(printf -- "node_index cohort_name cohort_id query_id\n0 testCohort -1 -1\n1 testCohort -1 -1\n2 testCohort -1 -1")"
getSavedCohort2="$(printf -- "node_index cohort_name query_id\n0 testCohort -1\n0 testCohort2 -1\n1 testCohort -1\n1 testCohort2 -1\n2 testCohort -1\n2 testCohort2 -1")"

# test2
timerHeaders="node_index,timer_description,duration_milliseconds"
survivalDays="$(printf -- "time_granularity,node_index,group_id,initial_count,time_point,event_of_interest_count,censoring_event_count\nday,0,Full cohort,684,0,0,0\nday,0,Full cohort,684,1,0,0\nday,0,Full cohort,684,2,0,0\nday,0,Full cohort,684,3,0,0\nday,0,Full cohort,684,4,0,0\nday,0,Full cohort,684,5,3,0")"
survivalWeeks="$(printf -- "time_granularity,node_index,group_id,initial_count,time_point,event_of_interest_count,censoring_event_count\nweek,0,Full cohort,684,0,0,0\nweek,0,Full cohort,684,1,3,0\nweek,0,Full cohort,684,2,18,0\nweek,0,Full cohort,684,3,3,0\nweek,0,Full cohort,684,4,3,0\nweek,0,Full cohort,684,5,6,0")"
survivalMonths="$(printf -- "time_granularity,node_index,group_id,initial_count,time_point,event_of_interest_count,censoring_event_count\nmonth,0,Full cohort,684,0,0,0\nmonth,0,Full cohort,684,1,30,0\nmonth,0,Full cohort,684,2,21,0\nmonth,0,Full cohort,684,3,30,0\nmonth,0,Full cohort,684,4,30,6\nmonth,0,Full cohort,684,5,30,0")"
survivalYears="$(printf -- "time_granularity,node_index,group_id,initial_count,time_point,event_of_interest_count,censoring_event_count\nyear,0,Full cohort,684,0,0,0\nyear,0,Full cohort,684,1,363,126\nyear,0,Full cohort,684,2,114,42\nyear,0,Full cohort,684,3,18,21\nyear,0,Full cohort,684,4,0,0\nyear,0,Full cohort,684,5,0,0")"

# test3
survivalSubGroup1="$(printf -- "time_granularity,node_index,group_id,initial_count,time_point,event_of_interest_count,censoring_event_count\nweek,0,Female,414,0,0,0\nweek,0,Female,414,1,0,0\nweek,0,Female,414,2,18,0\nweek,0,Female,414,3,3,0\nweek,0,Female,414,4,3,0\nweek,0,Female,414,5,6,0")"
survivalSubGroup2="$(printf -- "week,0,Male,270,0,0,0\nweek,0,Male,270,1,3,0\nweek,0,Male,270,2,0,0\nweek,0,Male,270,3,0,0\nweek,0,Male,270,4,0,0\nweek,0,Male,270,5,0,0")"

test1 () {
  docker-compose -f docker-compose.tools.yml run medco-cli-client --user $USERNAME --password $PASSWORD --o /results/result.csv get-saved-cohorts
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
  docker-compose -f docker-compose.tools.yml run medco-cli-client --user $USERNAME --password $PASSWORD --o /results/result.csv get-saved-cohorts
  result="$(awk -F',' '{print $1,$2,$4}' ../result.csv)"
  if [ "${result}" != "${getSavedCohort2}" ];
  then
  echo "get-saved-cohorts content after added new cohort: test failed"
  echo "result: ${result}" && echo "expected result: ${getSavedCohort2}"
  exit 1
  fi


  docker-compose -f docker-compose.tools.yml run medco-cli-client --user $USERNAME --password $PASSWORD update-saved-cohorts -c testCohort2 -p $(echo -1,-1,-1)
  docker-compose -f docker-compose.tools.yml run medco-cli-client --user $USERNAME --password $PASSWORD --o /results/result.csv get-saved-cohorts
  result="$(awk -F',' '{print $1,$2,$4}' ../result.csv)"
  if [ "${result}" != "${getSavedCohort2}" ];
  then
  echo "get-saved-cohorts content after update: test failed"
  echo "result: ${result}" && echo "expected result: ${getSavedCohort2}"
  exit 1
  fi

  docker-compose -f docker-compose.tools.yml run medco-cli-client --user $USERNAME --password $PASSWORD remove-saved-cohorts -c testCohort2
  docker-compose -f docker-compose.tools.yml run medco-cli-client --user $USERNAME --password $PASSWORD --o /results/result.csv get-saved-cohorts
  result="$(awk -F',' '{print $1,$2,$3,$4}' ../result.csv)"
  if [ "${result}" != "${getSavedCohort1}" ];
  then
  echo "get-saved-cohorts content after update: test failed"
  echo "result: ${result}" && echo "expected result: ${getSavedCohort1}"
  exit 1
  fi

}

test2 () {
  docker-compose -f docker-compose.tools.yml run medco-cli-client --user $USERNAME --password $PASSWORD -o /results/result.csv srva  -c testCohort -l 6 -g ${1}  -s /SPHN/SPHNv2020.1/FophDiagnosis/ -e /SPHN/SPHNv2020.1/DeathStatus/ -y 126:1 -d /results/timers.csv
  
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

test3 () {
  docker-compose -f docker-compose.tools.yml run \
    -v "../../test/survival_test_parameters.yaml":/parameters/survival_test_parameters.yaml \
    medco-cli-client --user $USERNAME --password $PASSWORD -o /results/result.csv srva -d /results/timers.csv \
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
echo "Testing saved-cohorts features..."

test1

test2 "day" "${survivalDays}"
test2 "week" "${survivalWeeks}"
test2 "month" "${survivalMonths}"
test2 "year" "${survivalYears}"

test3 "${survivalSubGroup1}" "${survivalSubGroup2}"

echo "CLI test 1/2 successful!"
popd
exit 0