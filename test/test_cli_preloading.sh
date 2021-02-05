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

# test4
getSavedCohortHeaders="node_index,cohort_name,cohort_id,query_id,creation_date,update_date,query_timing,panels"
getSavedCohort1="$(printf -- "node_index cohort_name cohort_id query_id query_timing panels\n\
0 testCohort -1 -1 any \"{panels:[{items:[{encrypted:false,queryTerm:/E2ETEST/SPHNv2020.1/DeathStatus/}],not:false,panelTiming:any}]}\"\n\
1 testCohort -1 -1 any \"{panels:[{items:[{encrypted:false,queryTerm:/E2ETEST/SPHNv2020.1/DeathStatus/}],not:false,panelTiming:any}]}\"\n\
2 testCohort -1 -1 any \"{panels:[{items:[{encrypted:false,queryTerm:/E2ETEST/SPHNv2020.1/DeathStatus/}],not:false,panelTiming:any}]}\"")"
getSavedCohort2="$(printf -- "node_index cohort_name query_id\n0 testCohort2 -1\n0 testCohort -1\n1 testCohort2 -1\n1 testCohort -1\n2 testCohort2 -1\n2 testCohort -1")"

# test5
timerHeaders="node_index,timer_description,duration_milliseconds"
survivalDays="$(printf -- "time_granularity,node_index,group_id,initial_count,time_point,event_of_interest_count,censoring_event_count\nday,0,Full cohort,684,0,0,0\nday,0,Full cohort,684,1,0,0\nday,0,Full cohort,684,2,0,0\nday,0,Full cohort,684,3,0,0\nday,0,Full cohort,684,4,0,0\nday,0,Full cohort,684,5,3,0")"
survivalWeeks="$(printf -- "time_granularity,node_index,group_id,initial_count,time_point,event_of_interest_count,censoring_event_count\nweek,0,Full cohort,684,0,0,0\nweek,0,Full cohort,684,1,3,0\nweek,0,Full cohort,684,2,18,0\nweek,0,Full cohort,684,3,3,0\nweek,0,Full cohort,684,4,3,0\nweek,0,Full cohort,684,5,6,0")"
survivalMonths="$(printf -- "time_granularity,node_index,group_id,initial_count,time_point,event_of_interest_count,censoring_event_count\nmonth,0,Full cohort,684,0,0,0\nmonth,0,Full cohort,684,1,30,0\nmonth,0,Full cohort,684,2,21,0\nmonth,0,Full cohort,684,3,30,0\nmonth,0,Full cohort,684,4,30,6\nmonth,0,Full cohort,684,5,30,0")"
survivalYears="$(printf -- "time_granularity,node_index,group_id,initial_count,time_point,event_of_interest_count,censoring_event_count\nyear,0,Full cohort,684,0,0,0\nyear,0,Full cohort,684,1,363,126\nyear,0,Full cohort,684,2,114,42\nyear,0,Full cohort,684,3,18,21\nyear,0,Full cohort,684,4,0,0\nyear,0,Full cohort,684,5,0,0")"

# test6
survivalSubGroup1="$(printf -- "time_granularity,node_index,group_id,initial_count,time_point,event_of_interest_count,censoring_event_count\nweek,0,Female,414,0,0,0\nweek,0,Female,414,1,0,0\nweek,0,Female,414,2,18,0\nweek,0,Female,414,3,3,0\nweek,0,Female,414,4,3,0\nweek,0,Female,414,5,6,0")"
survivalSubGroup2="$(printf -- "week,0,Male,270,0,0,0\nweek,0,Male,270,1,3,0\nweek,0,Male,270,2,0,0\nweek,0,Male,270,3,0,0\nweek,0,Male,270,4,0,0\nweek,0,Male,270,5,0,0")"

# test7
function timing() { echo "query clr::/E2ETEST/SPHNv2020.1/DeathStatus/ OR clr::/E2ETEST/SPHNv2020.1/DeathStatus/ ${1} AND clr::/E2ETEST/SPHNv2020.1/DeathStatus/:/E2ETEST/DeathStatus-status/death/:/SPHNv2020.1/DeathStatus/ ${2} AND clr::/E2ETEST/I2B2/Demographics/Gender/Female/ OR clr::/E2ETEST/I2B2/Demographics/Gender/Male/ ${3} -t ${4}"; };
timingResultNonZeroExpected="$(printf -- "count\n165\n165\n165")"
timingResultZeroExpected="$(printf -- "count\n0\n0\n0")"

# test8
function cohortPatientListWithCredentials() { docker-compose -f docker-compose.tools.yml run medco-cli-client --user ${1} --password ${2} --o /data/result.csv \
  cpl -c testCohort -d /data/timers.csv; };
patientList="$(printf -- "Node idx 0\n1137,1138,1139,1140,1141,1142,1143,1144,1145,1146,1147,1148,1149,1150,1151,1152,1153,1154,1155,1156,1157,1158,1159,1160,1161,1162,1163,1164,1165,1166,1167,1168,1169,1170,1171,1172,1173,1174,1175,1176,1177,1178,1179,1180,1181,1182,1183,1184,1185,1186,1187,1188,1189,1190,1191,1192,1193,1194,1195,1196,1197,\
1198,1199,1200,1201,1202,1203,1204,1205,1206,1207,1208,1209,1210,1211,1212,1213,1214,1215,1216,1217,1218,1219,1220,1221,1222,1223,1224,1225,1226,1227,1228,1229,1230,1231,1232,1233,1234,1235,1236,1237,1238,1239,1240,1241,1242,1243,1244,1245,1246,1247,1248,1249,1250,1251,1252,1253,1254,1255,1256,1257,1258,1259,1260,1261,1262,1263,1264,1265,1266,\
1267,1268,1269,1270,1271,1272,1273,1274,1275,1276,1277,1278,1279,1280,1281,1282,1283,1284,1285,1286,1287,1288,1289,1290,1291,1292,1293,1294,1295,1296,1297,1298,1299,1300,1301,1302,1303,1304,1305,1306,1307,1308,1309,1310,1311,1312,1313,1314,1315,1316,1317,1318,1319,1320,1321,1322,1323,1324,1325,1326,1327,1328,1329,1330,1331,1332,1333,1334,1335,\
1336,1337,1338,1339,1340,1341,1342,1343,1344,1345,1346,1347,1348,1349,1350,1351,1352,1353,1354,1355,1356,1357,1358,1359,1360,1361,1362,1363,1364")"
expectedError="is not authorized to query patient lists"

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

  result="$(awk -vFPAT='("[^"]+")|([^,]+)' '{print $1,$2,$3,$4,$7,$8}' ../result.csv)"
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
  result="$(awk -vFPAT='("[^"]+")|([^,]+)' '{print $1,$2,$3,$4,$7,$8}' ../result.csv)"
  if [ "${result}" != "${getSavedCohort1}" ];
  then
  echo "get-saved-cohorts content after removing new cohorts: test failed"
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
  echo "timer headers for survival $1: test failed"
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

test8() {

  cohortPatientListWithCredentials "test" "test"

  result="$(awk -F',' 'NR==1, NR==2 {print $0}' ../result.csv)"
  if [ "${result}" != "${patientList}" ]
  then
  echo "cohorts patient list: test failed"
  echo "result: ${result}" && echo "expected result: ${patientList}"
  exit 1
  fi

  result="$(awk -F',' 'NR==1{print $0}' ../timers.csv)"

  if [ "${result}" != "${timerHeaders}" ];
  then
  echo "timer headers for cohort patient list: test failed"
  echo "result: ${result}" && echo "expected result: ${timerHeaders}"
  exit 1
  fi

  # test when the application is not supposed to authorize the user to get patient lists
  non_authorized="test_explore_count_global"
  # empty prievious result file, create a dump file for logs
  echo "" > ../result.csv
  echo "" > ../dumped_logs_to_remove.txt


  cohortPatientListWithCredentials "${non_authorized}" "test" > ../dumped_logs_to_remove.txt 2>&1

  # result file must be empty
  if [[ $(cat ../result.csv)  != "" ]];
  then
  echo "cohorts patient list: test failed"
  echo "result file must remain empty for non authorized user ${non_authorized}"
  exit 1
  fi

  # check the error description
  description="$(awk 'END{print}' ../dumped_logs_to_remove.txt)"
  if [[ "${description}" != *"${expectedError}"* ]];
  then
  echo "cohorts patient list: test failed"
  echo "last line of log is ${description}: expected to contain  \"${expectedError}\""
  exit 1
  fi

  rm ../dumped_logs_to_remove.txt  
}

pushd deployments/dev-local-3nodes/
echo "Testing concept-children..."

test1 "concept-children" "${searchConceptChildren1}" "${resultSearchConceptChildren1}"
test1 "concept-children" "${searchConceptChildren2}" "${resultSearchConceptChildren2}"

echo "Testing modifier-children..."

test1 "modifier-children" "${searchModifierChildren}" "${resultSearchModifierChildren}"

echo "Testing concept-info..."

test1 "concept-info" "${searchConceptInfo}" "${resultSearchConceptInfo}"

echo "Testing modifier-info..."

test1 "modifier-info" "${searchModifierInfo}" "${resultSearchModifierInfo}"

echo "Testing query with test user..."

test2 "query " "${query1}" "${resultQuery1}"
test2 "query " "${query2}" "${resultQuery2}"
test2 "query " "${query3}" "${resultQuery3}"
test2 "query " "${query4}" "${resultQuery4}"
test2 "query " "${query5}" "${resultQuery5a}"
test2 "query " "${query6}" "${resultQuery6}"
test2 "query " "${query7}" "${resultQuery7}"
test2 "query " "${query8}" "${resultQuery8}"
test2 "query " "${query9}" "${resultQuery9}"
test2 "query " "${query10}" "${resultQuery10}"
test2 "query " "${query11}" "${resultQuery11}"
test2 "query " "${query12}" "${resultQuery12}"
test2 "query " "${query13}" "${resultQuery13}"

echo "Testing query with test_explore_patient_list user..."
USERNAME="${1:-test}_explore_patient_list"

test2 "query " "${query1}" "${resultQuery1}"
test2 "query " "${query2}" "${resultQuery2}"
test2 "query " "${query3}" "${resultQuery3}"
test2 "query " "${query4}" "${resultQuery4}"
test2 "query " "${query5}" "${resultQuery5a}"

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

echo "Testing cohorts-patient-list"

test8 

echo "CLI test 1/2 successful!"
popd
exit 0