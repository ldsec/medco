#!/usr/bin/env bash
set -Eeuo pipefail

USERNAME=${1:-test}
PASSWORD=${2:-test}

getSavedCohortHeaders="node_index,cohort_name,cohort_id,query_id,creation_date,update_date,query_timing,query_timing_sequence,panels"
getSavedCohort1="$(printf -- "node_index cohort_name cohort_id query_id query_timing query_timing_sequence panels\n\
0 testCohort -1 -1 any {temporalSequence:[]} \"{panels:[{cohortItems:null,conceptItems:[{encrypted:false,queryTerm:/E2ETEST/SPHNv2020.1/DeathStatus/}],not:false,panelTiming:any}]}\"\n\
1 testCohort -1 -1 any {temporalSequence:[]} \"{panels:[{cohortItems:null,conceptItems:[{encrypted:false,queryTerm:/E2ETEST/SPHNv2020.1/DeathStatus/}],not:false,panelTiming:any}]}\"\n\
2 testCohort -1 -1 any {temporalSequence:[]} \"{panels:[{cohortItems:null,conceptItems:[{encrypted:false,queryTerm:/E2ETEST/SPHNv2020.1/DeathStatus/}],not:false,panelTiming:any}]}\"")"
getSavedCohort2="$(printf -- "node_index cohort_name query_id\n0 testCohort2 -1\n0 testCohort -1\n1 testCohort2 -1\n1 testCohort -1\n2 testCohort2 -1\n2 testCohort -1")"
test1 () {
  docker-compose -f docker-compose.tools.yml run medco-cli-client --user $USERNAME --password $PASSWORD --o /data/result.csv get-saved-cohorts
  result="$(awk -F',' 'NR==1{print $0}' ../result.csv)"
  if [ "${result}" != "${getSavedCohortHeaders}" ];
  then
  echo "get-saved-cohorts headers: test failed"
  echo "result: ${result}" && echo "expected result: ${getSavedCohortHeaders}"
  exit 1
  fi

  result="$(awk -vFPAT='("[^"]+")|([^,]+)' '{print $1,$2,$3,$4,$7,$8,$9}' ../result.csv)"
  if [ "${result}" != "${getSavedCohort1}" ];
  then
  echo "get-saved-cohorts content before update: test failed"
  echo "result: ${result}" && echo "expected result: ${getSavedCohort1}"
  exit 1
  fi

  docker-compose -f docker-compose.tools.yml run medco-cli-client --user $USERNAME --password $PASSWORD add-saved-cohorts -c testCohort2 -q $(echo -1,-1,-1)
  docker-compose -f docker-compose.tools.yml run medco-cli-client --user $USERNAME --password $PASSWORD --o /data/result.csv get-saved-cohorts
  result="$(awk -F',' '{print $1,$2,$4}' ../result.csv)"
  if [ "${result}" != "${getSavedCohort2}" ];
  then
  echo "get-saved-cohorts content after added new cohort: test failed"
  echo "result: ${result}" && echo "expected result: ${getSavedCohort2}"
  exit 1
  fi


  docker-compose -f docker-compose.tools.yml run medco-cli-client --user $USERNAME --password $PASSWORD update-saved-cohorts -c testCohort2 -q $(echo -1,-1,-1)
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
  result="$(awk -vFPAT='("[^"]+")|([^,]+)' '{print $1,$2,$3,$4,$7,$8,$9}' ../result.csv)"
  if [ "${result}" != "${getSavedCohort1}" ];
  then
  echo "get-saved-cohorts content after removing new cohorts: test failed"
  echo "result: ${result}" && echo "expected result: ${getSavedCohort1}"
  exit 1
  fi

}

timerHeaders="node_index,timer_description,duration_milliseconds"
survivalDays="$(printf -- "time_granularity,node_index,group_id,initial_count,time_point,event_of_interest_count,censoring_event_count\nday,0,Full cohort,684,0,0,0\nday,0,Full cohort,684,1,0,0\nday,0,Full cohort,684,2,0,0\nday,0,Full cohort,684,3,0,0\nday,0,Full cohort,684,4,0,0\nday,0,Full cohort,684,5,3,0")"
survivalWeeks="$(printf -- "time_granularity,node_index,group_id,initial_count,time_point,event_of_interest_count,censoring_event_count\nweek,0,Full cohort,684,0,0,0\nweek,0,Full cohort,684,1,3,0\nweek,0,Full cohort,684,2,18,0\nweek,0,Full cohort,684,3,3,0\nweek,0,Full cohort,684,4,3,0\nweek,0,Full cohort,684,5,6,0")"
survivalMonths="$(printf -- "time_granularity,node_index,group_id,initial_count,time_point,event_of_interest_count,censoring_event_count\nmonth,0,Full cohort,684,0,0,0\nmonth,0,Full cohort,684,1,30,0\nmonth,0,Full cohort,684,2,21,0\nmonth,0,Full cohort,684,3,30,0\nmonth,0,Full cohort,684,4,30,6\nmonth,0,Full cohort,684,5,30,0")"
survivalYears="$(printf -- "time_granularity,node_index,group_id,initial_count,time_point,event_of_interest_count,censoring_event_count\nyear,0,Full cohort,684,0,0,0\nyear,0,Full cohort,684,1,363,126\nyear,0,Full cohort,684,2,114,42\nyear,0,Full cohort,684,3,18,21\nyear,0,Full cohort,684,4,0,0\nyear,0,Full cohort,684,5,0,0")"
test2 () {
  docker-compose -f docker-compose.tools.yml run medco-cli-client --user $USERNAME --password $PASSWORD -o /data/result.csv srva  -c testCohort -l 6 -g ${1} \
   -s clr::/SPHN/SPHNv2020.1/FophDiagnosis/ \
   -w first \
   -e clr::/SPHN/SPHNv2020.1/DeathStatus/:/SPHN/DeathStatus-status/death/:/SPHNv2020.1/DeathStatus/ \
   -z earliest \
   -d /data/timers.csv

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

survivalSubGroup1="$(printf -- "time_granularity,node_index,group_id,initial_count,time_point,event_of_interest_count,censoring_event_count\nweek,0,Female,414,0,0,0\nweek,0,Female,414,1,0,0\nweek,0,Female,414,2,18,0\nweek,0,Female,414,3,3,0\nweek,0,Female,414,4,3,0\nweek,0,Female,414,5,6,0")"
survivalSubGroup2="$(printf -- "week,0,Male,270,0,0,0\nweek,0,Male,270,1,3,0\nweek,0,Male,270,2,0,0\nweek,0,Male,270,3,0,0\nweek,0,Male,270,4,0,0\nweek,0,Male,270,5,0,0")"
test3 () {
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

survivalSequenceOfEvents="$(printf -- "All 495\nNone 0")"
test4 () {
  docker-compose -f docker-compose.tools.yml run \
    -v "${PWD}/../../test/survival_e2e_test_parameters_sequence_of_events.yaml":/parameters/survival_e2e_test_parameters.yaml \
    medco-cli-client --user $USERNAME --password $PASSWORD -o /data/result.csv srva -d /data/timers.csv \
    -p /parameters/survival_e2e_test_parameters.yaml

    result="$(awk -F',' 'NR==7, NR==8 {print $3, $4}' ../result.csv)"
    if [ "${result}" != "${1}" ];
  then
  echo "survival analysis sequence of events in sub groups: test failed"
  echo "result: ${result}" && echo "expected result: ${1}"
  exit 1
  fi

}

function cohortPatientListWithCredentials() { docker-compose -f docker-compose.tools.yml run medco-cli-client --user ${1} --password ${2} --o /data/result.csv \
  cpl -c testCohort -d /data/timers.csv; };
patientList="$(printf -- "Node idx 0\n1137,1138,1139,1140,1141,1142,1143,1144,1145,1146,1147,1148,1149,1150,1151,1152,1153,1154,1155,1156,1157,1158,1159,1160,1161,1162,1163,1164,1165,1166,1167,1168,1169,1170,1171,1172,1173,1174,1175,1176,1177,1178,1179,1180,1181,1182,1183,1184,1185,1186,1187,1188,1189,1190,1191,1192,1193,1194,1195,1196,1197,\
1198,1199,1200,1201,1202,1203,1204,1205,1206,1207,1208,1209,1210,1211,1212,1213,1214,1215,1216,1217,1218,1219,1220,1221,1222,1223,1224,1225,1226,1227,1228,1229,1230,1231,1232,1233,1234,1235,1236,1237,1238,1239,1240,1241,1242,1243,1244,1245,1246,1247,1248,1249,1250,1251,1252,1253,1254,1255,1256,1257,1258,1259,1260,1261,1262,1263,1264,1265,1266,\
1267,1268,1269,1270,1271,1272,1273,1274,1275,1276,1277,1278,1279,1280,1281,1282,1283,1284,1285,1286,1287,1288,1289,1290,1291,1292,1293,1294,1295,1296,1297,1298,1299,1300,1301,1302,1303,1304,1305,1306,1307,1308,1309,1310,1311,1312,1313,1314,1315,1316,1317,1318,1319,1320,1321,1322,1323,1324,1325,1326,1327,1328,1329,1330,1331,1332,1333,1334,1335,\
1336,1337,1338,1339,1340,1341,1342,1343,1344,1345,1346,1347,1348,1349,1350,1351,1352,1353,1354,1355,1356,1357,1358,1359,1360,1361,1362,1363,1364")"
expectedError="is not authorized to query patient lists"
test5() {

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
echo "Testing saved-cohorts features..."
test1

echo "Testing survival analysis features..."
test2 "day" "${survivalDays}"
test2 "week" "${survivalWeeks}"
test2 "month" "${survivalMonths}"
test2 "year" "${survivalYears}"

test3 "${survivalSubGroup1}" "${survivalSubGroup2}"

test4 "${survivalSequenceOfEvents}"

echo "Testing cohorts-patient-list"
test5

popd
exit 0