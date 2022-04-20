#!/usr/bin/env bash
set -Eeuo pipefail

psql $PSQL_PARAMS -d "$MC_DB_NAME" <<-EOSQL

INSERT INTO query_tools.explore_query_results(
        query_id, query_name, user_id, query_status, clear_result_set_size, clear_result_set, query_definition, i2b2_encrypted_patient_set_id, i2b2_non_encrypted_patient_set_id)
        VALUES (-1, 'SurvivalE2ETest', 'test', 'completed', 228,'{1137,1138,1139,1140,1141,1142,1143,1144,1145,1146,1147,1148,1149,1150,1151,1152,1153,1154,1155,1156,1157,1158,1159,1160,1161,1162,1163,1164,1165,1166,1167,1168,1169,1170,1171,1172,1173,1174,1175,1176,1177,1178,1179,1180,1181,1182,1183,1184,1185,1186,1187,1188,1189,1190,1191,1192,1193,1194,1195,1196,1197,1198,1199,1200,1201,1202,1203,1204,1205,1206,1207,1208,1209,1210,1211,1212,1213,1214,1215,1216,1217,1218,1219,1220,1221,1222,1223,1224,1225,1226,1227,1228,1229,1230,1231,1232,1233,1234,1235,1236,1237,1238,1239,1240,1241,1242,1243,1244,1245,1246,1247,1248,1249,1250,1251,1252,1253,1254,1255,1256,1257,1258,1259,1260,1261,1262,1263,1264,1265,1266,1267,1268,1269,1270,1271,1272,1273,1274,1275,1276,1277,1278,1279,1280,1281,1282,1283,1284,1285,1286,1287,1288,1289,1290,1291,1292,1293,1294,1295,1296,1297,1298,1299,1300,1301,1302,1303,1304,1305,1306,1307,1308,1309,1310,1311,1312,1313,1314,1315,1316,1317,1318,1319,1320,1321,1322,1323,1324,1325,1326,1327,1328,1329,1330,1331,1332,1333,1334,1335,1336,1337,1338,1339,1340,1341,1342,1343,1344,1345,1346,1347,1348,1349,1350,1351,1352,1353,1354,1355,1356,1357,1358,1359,1360,1361,1362,1363,1364}', '{"panels":[{"cohortItems":null,"conceptItems":[{"encrypted":false,"queryTerm":"/E2ETEST/SPHNv2020.1/DeathStatus/"}],"not":false,"panelTiming":"any"}],"queryTiming":"any"}', -1, -1)
        ON CONFLICT DO NOTHING;
EOSQL




#Load test cohort

psql $PSQL_PARAMS -d "$MC_DB_NAME" <<-EOSQL

INSERT INTO query_tools.saved_cohorts(
        cohort_id, user_id, cohort_name, query_id, create_date, update_date, predefined, default_flag)
        VALUES (-1, 'test', 'testCohort', -1, '2020-08-25 13:57:00', '2020-08-25 13:57:00', FALSE, FALSE)
        ON CONFLICT DO NOTHING;

EOSQL