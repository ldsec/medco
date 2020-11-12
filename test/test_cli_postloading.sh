#!/usr/bin/env bash
set -Eeuo pipefail

USERNAME=${1:-test}
PASSWORD=${2:-test}

# test1
query1="enc::1 AND enc::2"
resultQuery1="$(printf -- "count\n8\n8\n8")"
query2="enc::6 OR enc::16 AND enc::8"
resultQuery2="$(printf -- "count\n5\n5\n5")"
query3="enc::5 AND enc::10 AND enc::15"
resultQuery3="$(printf -- "count\n3\n3\n3")"
query4="enc::4 OR enc::11 OR enc::17"
resultQuery4="$(printf -- "count\n7\n7\n7")"
query5="enc::3 OR enc::6 AND enc::9 AND enc::12 OR enc::15"
resultQuery5="$(printf -- "count\n2\n2\n2")"

# test2
variantNameGetValuesValue="5238"
variantNameGetValuesResult="16:75238144:C>C
                            6:52380882:G>G"
proteinChangeGetValuesValue="g32"
proteinChangeGetValuesResult="G325R
                              G32E"
proteinChangeGetValuesValue2="7cfs*"
proteinChangeGetValuesResult2="S137Cfs*28"
hugoGeneSymbolGetValuesValue="tr5"
hugoGeneSymbolGetValuesResult="HTR5A"

variantNameGetVariantsValue="16:75238144:C>C"
variantNameGetVariantsResult="-4530899676219565056"
proteinChangeGetVariantsValue="G325R"
proteinChangeGetVariantsResult="-2429151887266669568"
hugoGeneSymbolGetVariantsValue="HTR5A"
hugoGeneSymbolGetVariantsResult1="-7039476204566471680
                                  -7039476580443220992
                                  -7039476780159200256"
hugoGeneSymbolGetVariantsResult2="-7039476204566471680
                                  -7039476580443220992"
hugoGeneSymbolGetVariantsResult3="-7039476780159200256"

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
  result="$(docker-compose -f docker-compose.tools.yml run -e LOG_LEVEL=1 -e CONN_TIMEOUT=10m medco-cli-client --user $USERNAME --password $PASSWORD $1 $2 | tr -d '\r\n\t ')"
  expectedResult="$(echo "${3}" | tr -d '\r\n\t ')"
  if [ "${result}" != "${expectedResult}" ];
  then
  echo "$1 $2: test failed"
  echo "result: ${result}" && echo "expected result: ${3}"
  exit 1
  fi
}

pushd deployments/dev-local-3nodes/

echo "Testing query with genomic data..."
USERNAME="${USERNAME}_explore_patient_list"

test1 "query " "${query1}" "${resultQuery1}"
test1 "query " "${query2}" "${resultQuery2}"
test1 "query " "${query3}" "${resultQuery3}"
test1 "query " "${query4}" "${resultQuery4}"
test1 "query " "${query5}" "${resultQuery5}"

echo "Testing ga-get-values..."
USERNAME=${1:-test}

test2 "ga-get-values variant_name" "${variantNameGetValuesValue}" "${variantNameGetValuesResult}"
test2 "ga-get-values protein_change" "${proteinChangeGetValuesValue}" "${proteinChangeGetValuesResult}"
test2 "ga-get-values protein_change" "${proteinChangeGetValuesValue2}" "${proteinChangeGetValuesResult2}"
test2 "ga-get-values hugo_gene_symbol" "${hugoGeneSymbolGetValuesValue}" "${hugoGeneSymbolGetValuesResult}"

echo "Testing ga-get-variant..."

test2 "ga-get-variant variant_name" "${variantNameGetVariantsValue}" "${variantNameGetVariantsResult}"
test2 "ga-get-variant protein_change" "${proteinChangeGetVariantsValue}" "${proteinChangeGetVariantsResult}"
test2 "ga-get-variant hugo_gene_symbol" "${hugoGeneSymbolGetVariantsValue}" "${hugoGeneSymbolGetVariantsResult1}"
test2 "ga-get-variant --z "heterozygous" hugo_gene_symbol" "${hugoGeneSymbolGetVariantsValue}" "${hugoGeneSymbolGetVariantsResult2}"
test2 "ga-get-variant --z "unknown" hugo_gene_symbol" "${hugoGeneSymbolGetVariantsValue}" "${hugoGeneSymbolGetVariantsResult3}"
test2 "ga-get-variant --z "heterozygous\|homozygous" hugo_gene_symbol" "${hugoGeneSymbolGetVariantsValue}" "${hugoGeneSymbolGetVariantsResult2}"
test2 "ga-get-variant --z "heterozygous\|unknown" hugo_gene_symbol" "${hugoGeneSymbolGetVariantsValue}" "${hugoGeneSymbolGetVariantsResult1}"
test2 "ga-get-variant --z "homozygous\|unknown" hugo_gene_symbol" "${hugoGeneSymbolGetVariantsValue}" "${hugoGeneSymbolGetVariantsResult3}"
test2 "ga-get-variant --z "heterozygous\|homozygous\|unknown" hugo_gene_symbol" "${hugoGeneSymbolGetVariantsValue}" "${hugoGeneSymbolGetVariantsResult1}"

echo "CLI test 2/2 successful!"
popd
exit 0