#!/usr/bin/env bash
set -Eeuo pipefail

USERNAME=${1:-test}
PASSWORD=${2:-test}


# test1
#If one of the tables doesn't appear in the result. Be sure it is present in the database. And if it is, be sure
#that the c_visualattribute column, related to that table, inside the table_access table is set to 'CA' and not 'CH' which hides the table when fetching the children of "/"
searchConceptChildren1="/"
resultSearchConceptChildren1="PATH  TYPE
                              /E2ETEST/e2etest/ concept_container 0
                              /I2B2/I2B2/ concept_container 0
                              /SPHN/SPHNv2020.1/  concept_container 0"

searchConceptChildren2="/E2ETEST/e2etest/"
resultSearchConceptChildren2="PATH  TYPE
                              /E2ETEST/e2etest/1/ concept 12
                              /E2ETEST/e2etest/2/ concept 12
                              /E2ETEST/e2etest/3/ concept 12
                              /E2ETEST/modifiers/ modifier_folder 12"

searchModifierChildren="/E2ETEST/modifiers/ /e2etest/% /E2ETEST/e2etest/1/"
resultSearchModifierChildren="PATH  TYPE
                              /E2ETEST/modifiers/1/ modifier 6"

searchConceptInfo="/E2ETEST/e2etest/1/"
resultSearchConceptInfo="  <ExploreSearchResultElement>
      <AppliedPath>@</AppliedPath>
      <Code>ENC_ID:1</Code>
      <DisplayName>E2E Concept 1</DisplayName>
      <Leaf>true</Leaf>
      <MedcoEncryption>
          <Encrypted>true</Encrypted>
          <ID>1</ID>
      </MedcoEncryption>
      <Metadata>
          <ValueMetadata>
              <ChildrenEncryptIDs></ChildrenEncryptIDs>
              <CreationDateTime></CreationDateTime>
              <DataType></DataType>
              <EncryptedType></EncryptedType>
              <EnumValues></EnumValues>
              <Flagstouse></Flagstouse>
              <NodeEncryptID></NodeEncryptID>
              <Oktousevalues></Oktousevalues>
              <TestID></TestID>
              <TestName></TestName>
              <Version></Version>
          </ValueMetadata>
      </Metadata>
      <Name>E2E Concept 1</Name>
      <Path>/E2ETEST/e2etest/1/</Path>
      <Type>concept</Type>
  </ExploreSearchResultElement>"

searchModifierInfo="/E2ETEST/modifiers/1/ /e2etest/1/"
resultSearchModifierInfo="<ExploreSearchResultElement>
      <AppliedPath>/e2etest/1/</AppliedPath>
      <Code>ENC_ID:5</Code>
      <DisplayName>E2E Modifier 1</DisplayName>
      <Leaf>true</Leaf>
      <MedcoEncryption>
          <Encrypted>true</Encrypted>
          <ID>5</ID>
      </MedcoEncryption>
      <Metadata>
          <ValueMetadata>
              <ChildrenEncryptIDs></ChildrenEncryptIDs>
              <CreationDateTime></CreationDateTime>
              <DataType></DataType>
              <EncryptedType></EncryptedType>
              <EnumValues></EnumValues>
              <Flagstouse></Flagstouse>
              <NodeEncryptID></NodeEncryptID>
              <Oktousevalues></Oktousevalues>
              <TestID></TestID>
              <TestName></TestName>
              <Version></Version>
          </ValueMetadata>
      </Metadata>
      <Name>E2E Modifier 1</Name>
      <Path>/E2ETEST/modifiers/1/</Path>
      <Type>modifier</Type>
  </ExploreSearchResultElement>"

# test 2
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

query6="clr::/E2ETEST/e2etest/1/::EQ:10"
resultQuery6="$(printf -- "count\n1\n1\n1")"

query7="clr::/E2ETEST/e2etest/1/:/E2ETEST/modifiers/1/:/e2etest/1/::EQ:15"
resultQuery7="$(printf -- "count\n1\n1\n1")"

query8="clr::/E2ETEST/e2etest/1/::BETWEEN:5 and 25"
resultQuery8="$(printf -- "count\n2\n2\n2")"

query9="enc::1 OR clr::/E2ETEST/e2etest/2/::GE:25 AND clr::/E2ETEST/e2etest/2/:/E2ETEST/modifiers/2/:/e2etest/2/::LT:21"
resultQuery9="$(printf -- "count\n2\n2\n2")"

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

# test7
function timing() { echo "query clr::/E2ETEST/SPHNv2020.1/DeathStatus/ OR clr::/E2ETEST/SPHNv2020.1/DeathStatus/ ${1} AND clr::/E2ETEST/SPHNv2020.1/DeathStatus/:/E2ETEST/DeathStatus-status/death/:/SPHNv2020.1/DeathStatus/ ${2} AND clr::/E2ETEST/I2B2/Demographics/Gender/Female/ OR clr::/E2ETEST/I2B2/Demographics/Gender/Male/ ${3} -t ${4}"; };
timingResultNonZeroExpected="$(printf -- "count\n165\n165\n165")"
timingResultZeroExpected="$(printf -- "count\n0\n0\n0")"


#test whether each line from expected result is contained within the result.csv file
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
  docker-compose -f docker-compose.tools.yml run medco-cli-client --user $USERNAME --password $PASSWORD --o /data/result.csv $1 $2
  result="$(awk -F "\"*,\"*" '{print $2}' ../result.csv)"
  if [ "${result}" == "${3}" ];
  then
  echo "$1 $2: WARNING - result is the same"
  echo "result: ${result}" && echo "expected result: ${3}"
  fi
}

test4 () {
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

test5 () {
  docker-compose -f docker-compose.tools.yml run medco-cli-client --user $USERNAME --password $PASSWORD -o /data/result.csv srva  -c testCohort -l 6 -g ${1} \
   -s clr::/SPHN/SPHNv2020.1/FophDiagnosis/ \
   -e clr::/SPHN/SPHNv2020.1/DeathStatus/:/SPHN/DeathStatus-status/death/:/SPHNv2020.1/DeathStatus/ -d /data/timers.csv

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

test6 () {
  docker-compose -f docker-compose.tools.yml run \
    -v "${PWD}/../../test/survival_e2e_test_parameters.yaml":/parameters/survival_e2e_test_parameters.yaml \
    medco-cli-client --user $USERNAME --password $PASSWORD -o /data/result.csv srva -d /data/timers.csv \
    -p /parameters/survival_e2e_test_parameters.yaml

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

test7() {
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

pushd deployments/dev-local-3nodes/
printf "\nTesting concept-children...\n"

test1 "concept-children" "${searchConceptChildren1}" "${resultSearchConceptChildren1}"
test1 "concept-children" "${searchConceptChildren2}" "${resultSearchConceptChildren2}"

printf "\nTesting modifier-children...\n"

test1 "modifier-children" "${searchModifierChildren}" "${resultSearchModifierChildren}"

printf "\nTesting concept-info...\n"

test1 "concept-info" "${searchConceptInfo}" "${resultSearchConceptInfo}"

printf "\nTesting modifier-info...\n"

test1 "modifier-info" "${searchModifierInfo}" "${resultSearchModifierInfo}"

printf "\nTesting query with test user...\n"

test2 "query " "${query1}" "${resultQuery1}"
test2 "query " "${query2}" "${resultQuery2}"
test2 "query " "${query3}" "${resultQuery3}"
test2 "query " "${query4}" "${resultQuery4}"
test2 "query " "${query5}" "${resultQuery5a}"

echo "Testing query with test_explore_patient_list user..."
USERNAME="${1:-test}_explore_patient_list"

test2 "query " "${query1}" "${resultQuery1}"
test2 "query " "${query2}" "${resultQuery2}"
test2 "query " "${query3}" "${resultQuery3}"
test2 "query " "${query4}" "${resultQuery4}"
test2 "query " "${query5}" "${resultQuery5a}"
test2 "query " "${query6}" "${resultQuery6}"
test2 "query " "${query7}" "${resultQuery7}"
test2 "query " "${query8}" "${resultQuery8}"
test2 "query " "${query9}" "${resultQuery9}"

USERNAME="${1:-test}_explore_count_global"
test2 "query " "${query5}" "${resultQuery5b}"

USERNAME="${1:-test}_explore_count_global_obfuscated"
test3 "query " "${query5}" "${resultQuery5b}"

echo "Testing saved-cohorts features..."
USERNAME=${1:-test}

test4

echo "Testing survival analysis features..."
USERNAME=${1:-test}
test5 "day" "${survivalDays}"
test5 "week" "${survivalWeeks}"
test5 "month" "${survivalMonths}"
test5 "year" "${survivalYears}"

test6 "${survivalSubGroup1}" "${survivalSubGroup2}"

echo "Testing query with timing settings features..."
USERNAME=${1:-test}
test7 "any" "any" "any" "any" "${timingResultNonZeroExpected}"
test7 "sameinstancenum" "sameinstancenum" "sameinstancenum" "sameinstancenum" "${timingResultZeroExpected}"
test7 "samevisit" "samevisit" "samevisit" "samevisit" "${timingResultZeroExpected}"
test7 "sameinstancenum" "sameinstancenum" "any" "sameinstancenum" "${timingResultNonZeroExpected}"
test7 "samevisit" "samevisit" "any" "samevisit" "${timingResultNonZeroExpected}"

echo "CLI test 1/2 successful!"
popd
exit 0