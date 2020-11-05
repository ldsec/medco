export PSQL_PARAMS="-h localhost -p 5432 -U i2b2"
export I2B2_DB_USER="i2b2";

for i in 0 1 2; do
	export I2B2_DB_NAME="i2b2medcosrv${i}" 

	psql $PSQL_PARAMS -d "$I2B2_DB_NAME" <<-EOSQL
		
		DELETE FROM medco_ont.table_access;

		--taken from 65-medco-ontology.sh 
		ALTER TABLE medco_ont.table_access OWNER TO $I2B2_DB_USER;

		insert into medco_ont.table_access (c_table_cd, c_table_name, c_protected_access, c_hlevel, c_fullname, c_name,
			c_synonym_cd, c_visualattributes, c_facttablecolumn, c_dimtablename,
			c_columnname, c_columndatatype, c_operator, c_dimcode, c_tooltip) VALUES
			('SENSITIVE_TAGGED', 'SENSITIVE_TAGGED', 'N', 1, '\medco\tagged\', 'MedCo Sensitive Tagged Ontology',
			'N', 'CH', 'concept_cd', 'concept_dimension', 'concept_path', 'T', 'LIKE', '\medco\tagged\', 'MedCo Sensitive Tagged Ontology');
		
		--from e2etest-survival.sh
		insert into medco_ont.table_access (c_table_cd, c_table_name, c_protected_access, c_hlevel, c_fullname, c_name,
			c_synonym_cd, c_visualattributes, c_facttablecolumn, c_dimtablename,
			c_columnname, c_columndatatype, c_operator, c_dimcode, c_tooltip) VALUES
			('SPHN', 'E2ETEST', 'N', '1', '\SPHNv2020.1\', 'SPHN Ontology version 2020.1',
			'N', 'CA', 'concept_cd', 'concept_dimension', 'concept_path', 'T', 'LIKE', '\e2etest\', 'SPHN Ontology version 2020.1'),
			('I2B2', 'E2ETEST', 'N', '1', '\I2B2\', 'I2B2 ontology',
			'N', 'CA', 'concept_cd', 'concept_dimension', 'concept_path', 'T', 'LIKE', '\e2etest\', 'I2B2 Ontology');

	
		DELETE FROM medco_ont.e2etest;

		DELETE FROM medco_ont.sensitive_tagged;

		DELETE FROM i2b2demodata_i2b2.concept_dimension;

		DELETE FROM i2b2demodata_i2b2.modifier_dimension;

		DELETE FROM i2b2demodata_i2b2.provider_dimension;

		DELETE FROM i2b2demodata_i2b2.patient_dimension;

		DELETE FROM i2b2demodata_i2b2.patient_mapping;

		DELETE FROM i2b2demodata_i2b2.visit_dimension;

		DELETE FROM i2b2demodata_i2b2.encounter_mapping;

		DELETE FROM i2b2demodata_i2b2.observation_fact;
	EOSQL

	#to execute in go/src/medco/
	bash build/package/i2b2/sql/80-e2etest-data.sh
	bash build/package/i2b2/sql/82-e2etest-survival.sh


	export I2B2_SQL_DIR="build/package/i2b2/sql"
	bash build/package/i2b2/sql/85-data-manipulation-functions.sh
	echo "done with manipulation functions"
done
