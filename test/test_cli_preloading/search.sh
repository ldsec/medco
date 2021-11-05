#!/usr/bin/env bash
set -Eeuo pipefail

USERNAME=${1:-test}
PASSWORD=${2:-test}

test () {
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

search1="foph"
resultSearch1="  <ExploreSearchResultElement>
      <AppliedPath>@</AppliedPath>
      <Code>A168</Code>
      <Comment></Comment>
      <DisplayName>Foph Diagnosis</DisplayName>
      <Leaf>true</Leaf>
      <MedcoEncryption>
          <Encrypted>false</Encrypted>
      </MedcoEncryption>
      <Metadata></Metadata>
      <Name>Foph Diagnosis</Name>
      <Parent>
          <AppliedPath>@</AppliedPath>
          <Code></Code>
          <Comment></Comment>
          <DisplayName>SPHN ontology</DisplayName>
          <Leaf>false</Leaf>
          <MedcoEncryption>
              <Encrypted>false</Encrypted>
          </MedcoEncryption>
          <Metadata></Metadata>
          <Name>SPHN ontology</Name>
          <Path>/SPHN/SPHNv2020.1/</Path>
          <Type>concept_container</Type>
      </Parent>
      <Path>/SPHN/SPHNv2020.1/FophDiagnosis/</Path>
      <Type>concept</Type>
  </ExploreSearchResultElement>"

search2="10"
resultSearch2="  <ExploreSearchResultElement>
      <AppliedPath>/SPHNv2020.1/FophDiagnosis/</AppliedPath>
      <Code>101:ICD10</Code>
      <Comment></Comment>
      <DisplayName>ICD10</DisplayName>
      <Leaf>false</Leaf>
      <MedcoEncryption>
          <Encrypted>false</Encrypted>
      </MedcoEncryption>
      <Metadata></Metadata>
      <Name>ICD10</Name>
      <Parent>
          <AppliedPath>/SPHNv2020.1/FophDiagnosis/</AppliedPath>
          <Code>101</Code>
          <Comment></Comment>
          <DisplayName>Diagnosis value</DisplayName>
          <Leaf>false</Leaf>
          <MedcoEncryption>
              <Encrypted>false</Encrypted>
          </MedcoEncryption>
          <Metadata></Metadata>
          <Name>Diagnosis value</Name>
          <Parent>
              <AppliedPath>@</AppliedPath>
              <Code>A168</Code>
              <Comment></Comment>
              <DisplayName>Foph Diagnosis</DisplayName>
              <Leaf>true</Leaf>
              <MedcoEncryption>
                  <Encrypted>false</Encrypted>
              </MedcoEncryption>
              <Metadata></Metadata>
              <Name>Foph Diagnosis</Name>
              <Parent>
                  <AppliedPath>@</AppliedPath>
                  <Code></Code>
                  <Comment></Comment>
                  <DisplayName>SPHN ontology</DisplayName>
                  <Leaf>false</Leaf>
                  <MedcoEncryption>
                      <Encrypted>false</Encrypted>
                  </MedcoEncryption>
                  <Metadata></Metadata>
                  <Name>SPHN ontology</Name>
                  <Path>/SPHN/SPHNv2020.1/</Path>
                  <Type>concept_container</Type>
              </Parent>
              <Path>/SPHN/SPHNv2020.1/FophDiagnosis/</Path>
              <Type>concept</Type>
          </Parent>
          <Path>/SPHN/FophDiagnosis-code/</Path>
          <Type>modifier_folder</Type>
      </Parent>
      <Path>/SPHN/FophDiagnosis-code/ICD10/</Path>
      <Type>modifier_folder</Type>
  </ExploreSearchResultElement>"

search3="gender"
resultSearch3="  <ExploreSearchResultElement>
      <AppliedPath>@</AppliedPath>
      <Code></Code>
      <Comment></Comment>
      <DisplayName>Gender</DisplayName>
      <Leaf>false</Leaf>
      <MedcoEncryption>
          <Encrypted>false</Encrypted>
      </MedcoEncryption>
      <Metadata></Metadata>
      <Name>Gender</Name>
      <Parent>
          <AppliedPath>@</AppliedPath>
          <Code></Code>
          <Comment></Comment>
          <DisplayName>I2B2 demographics</DisplayName>
          <Leaf>false</Leaf>
          <MedcoEncryption>
              <Encrypted>false</Encrypted>
          </MedcoEncryption>
          <Metadata></Metadata>
          <Name>I2B2 demographics</Name>
          <Parent>
              <AppliedPath>@</AppliedPath>
              <Code></Code>
              <Comment></Comment>
              <DisplayName>I2B2 ontology</DisplayName>
              <Leaf>false</Leaf>
              <MedcoEncryption>
                  <Encrypted>false</Encrypted>
              </MedcoEncryption>
              <Metadata></Metadata>
              <Name>I2B2 ontology</Name>
              <Path>/I2B2/I2B2/</Path>
              <Type>concept_container</Type>
          </Parent>
          <Path>/I2B2/I2B2/Demographics/</Path>
          <Type>concept_folder</Type>
      </Parent>
      <Path>/I2B2/I2B2/Demographics/Gender/</Path>
      <Type>concept_folder</Type>
  </ExploreSearchResultElement>
  <ExploreSearchResultElement>
      <AppliedPath>@</AppliedPath>
      <Code>DEM|SEX:f</Code>
      <Comment></Comment>
      <DisplayName>Female gender</DisplayName>
      <Leaf>true</Leaf>
      <MedcoEncryption>
          <Encrypted>false</Encrypted>
      </MedcoEncryption>
      <Metadata></Metadata>
      <Name>Female gender</Name>
      <Parent>
          <AppliedPath>@</AppliedPath>
          <Code></Code>
          <Comment></Comment>
          <DisplayName>Gender</DisplayName>
          <Leaf>false</Leaf>
          <MedcoEncryption>
              <Encrypted>false</Encrypted>
          </MedcoEncryption>
          <Metadata></Metadata>
          <Name>Gender</Name>
          <Parent>
              <AppliedPath>@</AppliedPath>
              <Code></Code>
              <Comment></Comment>
              <DisplayName>I2B2 demographics</DisplayName>
              <Leaf>false</Leaf>
              <MedcoEncryption>
                  <Encrypted>false</Encrypted>
              </MedcoEncryption>
              <Metadata></Metadata>
              <Name>I2B2 demographics</Name>
              <Parent>
                  <AppliedPath>@</AppliedPath>
                  <Code></Code>
                  <Comment></Comment>
                  <DisplayName>I2B2 ontology</DisplayName>
                  <Leaf>false</Leaf>
                  <MedcoEncryption>
                      <Encrypted>false</Encrypted>
                  </MedcoEncryption>
                  <Metadata></Metadata>
                  <Name>I2B2 ontology</Name>
                  <Path>/I2B2/I2B2/</Path>
                  <Type>concept_container</Type>
              </Parent>
              <Path>/I2B2/I2B2/Demographics/</Path>
              <Type>concept_folder</Type>
          </Parent>
          <Path>/I2B2/I2B2/Demographics/Gender/</Path>
          <Type>concept_folder</Type>
      </Parent>
      <Path>/I2B2/I2B2/Demographics/Gender/Female/</Path>
      <Type>concept</Type>
  </ExploreSearchResultElement>
  <ExploreSearchResultElement>
      <AppliedPath>@</AppliedPath>
      <Code>DEM|SEX:m</Code>
      <Comment></Comment>
      <DisplayName>Male gender</DisplayName>
      <Leaf>true</Leaf>
      <MedcoEncryption>
          <Encrypted>false</Encrypted>
      </MedcoEncryption>
      <Metadata></Metadata>
      <Name>Male gender</Name>
      <Parent>
          <AppliedPath>@</AppliedPath>
          <Code></Code>
          <Comment></Comment>
          <DisplayName>Gender</DisplayName>
          <Leaf>false</Leaf>
          <MedcoEncryption>
              <Encrypted>false</Encrypted>
          </MedcoEncryption>
          <Metadata></Metadata>
          <Name>Gender</Name>
          <Parent>
              <AppliedPath>@</AppliedPath>
              <Code></Code>
              <Comment></Comment>
              <DisplayName>I2B2 demographics</DisplayName>
              <Leaf>false</Leaf>
              <MedcoEncryption>
                  <Encrypted>false</Encrypted>
              </MedcoEncryption>
              <Metadata></Metadata>
              <Name>I2B2 demographics</Name>
              <Parent>
                  <AppliedPath>@</AppliedPath>
                  <Code></Code>
                  <Comment></Comment>
                  <DisplayName>I2B2 ontology</DisplayName>
                  <Leaf>false</Leaf>
                  <MedcoEncryption>
                      <Encrypted>false</Encrypted>
                  </MedcoEncryption>
                  <Metadata></Metadata>
                  <Name>I2B2 ontology</Name>
                  <Path>/I2B2/I2B2/</Path>
                  <Type>concept_container</Type>
              </Parent>
              <Path>/I2B2/I2B2/Demographics/</Path>
              <Type>concept_folder</Type>
          </Parent>
          <Path>/I2B2/I2B2/Demographics/Gender/</Path>
          <Type>concept_folder</Type>
      </Parent>
      <Path>/I2B2/I2B2/Demographics/Gender/Male/</Path>
      <Type>concept</Type>
  </ExploreSearchResultElement>"

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
      <Comment>E2E Concept 1</Comment>
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
      <Comment>E2E Modifier 1</Comment>
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

pushd deployments/dev-local-3nodes/
echo "Testing search..."

test "search" "${search1}" "${resultSearch1}"
test "search" "${search2}" "${resultSearch2}"
test "search" "${search3}" "${resultSearch3}"

echo "Testing concept-children..."

test "concept-children" "${searchConceptChildren1}" "${resultSearchConceptChildren1}"
test "concept-children" "${searchConceptChildren2}" "${resultSearchConceptChildren2}"

echo "Testing modifier-children..."

test "modifier-children" "${searchModifierChildren}" "${resultSearchModifierChildren}"

echo "Testing concept-info..."

test "concept-info" "${searchConceptInfo}" "${resultSearchConceptInfo}"

echo "Testing modifier-info..."

test "modifier-info" "${searchModifierInfo}" "${resultSearchModifierInfo}"

popd
exit 0
