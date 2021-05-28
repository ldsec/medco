"""

--looking for patients associated to concepts with 0 totalnum
SELECT c_fullname, patient_num
	FROM medco_ont.sphn, i2b2demodata_i2b2.observation_fact
	where c_totalnum is null 
		and c_basecode = concept_cd
	order by c_fullname limit 100 
# This could be a good unit test: i.e. checking that no rows are returned by the previous query.
# That is, there is no concept or modifier with a null totalnum that has any patient associated to it.
"""

# This is a script which is used to generate a postgresql script that will update the c_totalnum value of a concept or modifier
# so that c_totalnum contains the number of patients or observation (depends on your choice) linked directly or inderectly
# to the concept/modifier.

#tip: The generate psql script is easier to read than the postgreSQL_script that is waiting to be formatted by the execution of this script.




#TODO You need to treat c_facttablecolumn like it is case insensitive
postgreSQL_script = """
-- set client_min_messages to 'info';

 /* This script updates totalnum for all modifiers and concepts in the specified metadata table.
  *
  * The recursive formula for determining the totalnum value of a modifier is:
  * 	modifier totalnum = #(set of {subject}s directly referencing the modifier UNION of all sets of {subject}s referencing the children of the modifier )
  *
  * 	concept totalnum = #(set of {subject}s directly referencing the concept UNION of all sets of {subject}s referencing the children of the concept)
  *
  * Rules:
  * 1) A modifier is a parent of another modifier if the parent's c_fullname is the suffix of another child c_fullname.
  * 2) A concept is a parent of a modifier if the parent concept c_fullname is equal to the m_applied_path column of the child.
  * 3) Like rule 1 but the parent and child are both concepts
  *
  * Hence one may see multiple layers of modifiers and concepts
  *
  * The script uses a bottom approach to compute this formula
  */

/**
 * This code traverses iteratively the different heights of modfiers.
 * At each iteration the code adds the {subject}s that relate to the modifiers to the set of {subject}s that relate to the parent (concept or modifier).
 * After having done that totalnum is updated for the parent of the current iteration.
 */
CREATE OR REPLACE FUNCTION aggregateModifiersCounts() RETURNS void AS $$
	DECLARE
		min_depth int := (SELECT MIN(c_hlevel) FROM {metadata_schema_name}.{metadata_table_name} s WHERE s.c_facttablecolumn ILIKE 'MODIFIER_CD');
		max_depth int := (SELECT MAX(c_hlevel) FROM {metadata_schema_name}.{metadata_table_name} s WHERE s.c_facttablecolumn ILIKE 'MODIFIER_CD');
	BEGIN

		raise info 'before for loop';
		-- the children height is height and the parent height is height-1.
		FOR height IN REVERSE max_depth..min_depth
		LOOP
			raise info 'we are at height %', height;

			--link the {subject}s of the modifier children to their parent modifers or concept.
			INSERT INTO {metadata_schema_name}.concept_to_id
				SELECT DISTINCT cti.identifier, parent.c_fullname, (height-1)
				FROM {metadata_schema_name}.{metadata_table_name} child, {metadata_schema_name}.concept_to_id cti, {metadata_schema_name}.{metadata_table_name} parent
				WHERE cti.c_hlevel = height AND child.c_hlevel = height AND child.c_facttablecolumn ILIKE 'MODIFIER_CD'
					AND cti.c_fullname = child.c_fullname AND
					(
						( --the parent is a modifier at the height just above.
							parent.c_hlevel = (height - 1) AND parent.c_facttablecolumn ILIKE 'MODIFIER_CD' AND
							child.c_fullname LIKE (parent.c_fullname || '%') ESCAPE '|'
							-- Use | as an escape character instead of the backslash which appears in windows paths and breaks sql queries.
						)
						OR
						( --the parent is a concept
							parent.c_facttablecolumn ILIKE 'CONCEPT_CD' AND parent.c_fullname LIKE child.m_applied_path ESCAPE '|'
						)
					)
			;


		END LOOP;
	END;
$$ LANGUAGE plpgsql;

/**
 * This code traverses iteratively the different heights of concepts.
 * At each iteration the code adds the {subject}s that relate to the concepts to the set of {subject}s that relate to the parent concept.
 * After having done that totalnum is updated for the parents of the current iteration.
 */
CREATE OR REPLACE FUNCTION aggregateConceptsCounts() RETURNS void AS $$
	DECLARE
		min_depth int := (SELECT MIN(c_hlevel) FROM {metadata_schema_name}.{metadata_table_name} s WHERE s.c_facttablecolumn ILIKE 'CONCEPT_CD');
		max_depth int := (SELECT MAX(c_hlevel) FROM {metadata_schema_name}.{metadata_table_name} s WHERE s.c_facttablecolumn ILIKE 'CONCEPT_CD');
	BEGIN

		raise info 'before for loop';
		-- the children height is height and the parent height is height-1.
		FOR height IN REVERSE max_depth..min_depth
		LOOP
			raise info 'we are at height %', height;

			--link the {subject}s of the modifier children to their parent modifers or concept.
			INSERT INTO {metadata_schema_name}.concept_to_id
				SELECT DISTINCT cti.identifier, parent.c_fullname, (height-1)
				FROM {metadata_schema_name}.{metadata_table_name} child, {metadata_schema_name}.concept_to_id cti, {metadata_schema_name}.{metadata_table_name} parent
				WHERE cti.c_hlevel = height AND child.c_hlevel = height AND child.c_facttablecolumn ILIKE 'CONCEPT_CD'
					AND cti.c_fullname = child.c_fullname
					--the parent is a concept at the height just above.
					AND parent.c_hlevel = (height - 1) AND parent.c_facttablecolumn ILIKE 'CONCEPT_CD' AND
					child.c_fullname LIKE (parent.c_fullname || '%') ESCAPE '|'
					-- Use | as an escape character instead of the backslash which appears in windows paths and breaks sql queries.
			;

		END LOOP;
	END;
$$ LANGUAGE plpgsql;

DO $$
BEGIN

DROP TABLE IF EXISTS {metadata_schema_name}.concept_to_id;

--a temporary table used for computation of totalnum that links a modifer with a {subject}
CREATE TABLE {metadata_schema_name}.concept_to_id(
	identifier integer NOT NULL,
	c_fullname character varying(2000) COLLATE pg_catalog."default",
	c_hlevel integer NOT NULL
);


CREATE INDEX c_hlevel_cti
	ON {metadata_schema_name}.concept_to_id USING btree
	(c_hlevel ASC NULLS LAST)
	TABLESPACE pg_default;

CREATE INDEX c_fullname_idx_cti
	ON {metadata_schema_name}.concept_to_id USING btree
	(c_fullname COLLATE pg_catalog."default" ASC NULLS LAST)
	TABLESPACE pg_default;

ALTER TABLE {metadata_schema_name}.concept_to_id
	OWNER to {temporary_table_owner};

raise info 'created temporary table concept_to_id and its indexes';

-- filling totalnum for leaf modifiers.
UPDATE {metadata_schema_name}.{metadata_table_name} s
	SET c_totalnum = (
		SELECT COUNT (DISTINCT o.{observation_id} ) FROM {data_schema_name}.observation_fact o
		WHERE s.c_basecode = o.modifier_cd AND o.modifier_cd != '@'  -- '@' is for concepts
	)
WHERE s.c_facttablecolumn ILIKE 'MODIFIER_CD' AND s.c_basecode != '@';

raise info 'Done filling totalnum counting observations or patients directly linked to the modifiers';

--In this code we insert into concept_to_id the distinct (patient number, modfier fullname, modifier level) tuples that can be linked together by directly checking references between observation_fact and metadata
INSERT INTO {metadata_schema_name}.concept_to_id
SELECT DISTINCT o.{observation_id}, s.c_fullname, s.c_hlevel
FROM {data_schema_name}.observation_fact o, {metadata_schema_name}.{metadata_table_name} s
WHERE s.c_basecode = o.modifier_cd
	AND s.c_facttablecolumn ILIKE 'MODIFIER_CD' AND s.c_basecode != '@'
	AND o.modifier_cd != '@';  -- '@' is for concepts

--In this code we insert into concept_to_id the distinct (patient number, concept fullname, concept level) tuples that can be linked together by directly checking references between observation_fact and metadata
INSERT INTO {metadata_schema_name}.concept_to_id
SELECT DISTINCT o.{observation_id}, s.c_fullname, s.c_hlevel
FROM {data_schema_name}.observation_fact o, {metadata_schema_name}.{metadata_table_name} s
WHERE s.c_basecode = o.concept_cd AND o.modifier_cd = '@'  -- '@' is for concepts
	AND s.c_facttablecolumn ILIKE 'CONCEPT_CD';


PERFORM aggregateModifiersCounts();

raise info 'After call to aggregate modifiers {subject} function.';

PERFORM aggregateConceptsCounts();

raise info 'After call to aggregate concepts {subject} function.';

UPDATE {metadata_schema_name}.{metadata_table_name} s
SET c_totalnum = count_to_name.count -- count is either the count of distinct patients or the count of distinct observations per concept, modifier depending on what you what you chose to count upon.
FROM (
	SELECT COUNT(DISTINCT cti.identifier) AS count, cti.c_fullname AS fullname
	FROM {metadata_schema_name}.concept_to_id cti
	GROUP BY cti.c_fullname
) AS count_to_name
WHERE count_to_name.fullname = s.c_fullname;


--We are done with the temporary table.
-- DROP TABLE {metadata_schema_name}.concept_to_id;

--no need to explicitly commit since this is a implicit pgsql anonymous function which automatically commit the results at the END statement.
END;
$$ LANGUAGE plpgsql;

"""

print("""
This script will generate a postgreSQL script which you will be able to execute on your database in order to update
the totalnum column of your metadata/ontology table. We need a few informations about your database in order to generate that
SQL script:
_ The ontolgy schema name, that is the name of the schema that contains your metadata table.
_ The name of your metadata table which holds information about the modifiers and concepts such as the totalnum column.
_ The owner of the temporary table that will be used in order to perform the count query.
_ The data schema name, that is, the name of the schema which contains your observation_fact table.
_ Whether the c_facttablecolumn of the ontology contains upper-cased or lower-cased information (e.g. 'MODIFIER_CD' or 'modifier_cd').

You also have to choose which information the c_totalnum column of the ontology table will contain.
Option 1: The number of distinct patients associated to the concept/modifier.
Option 2: The number of distinct observation associated to the concept/modifier.
""")

import sys
if  len(sys.argv) > 1:
	metadata_schema_name = sys.argv[1]
	metadata_table_name = sys.argv[2]
	temporary_table_owner = sys.argv[3]
	data_schema_name = sys.argv[4]
	option = sys.argv[5]
else:
	while True:
		metadata_schema_name = input("metadata schema name =")
		metadata_table_name = input("metadata table name =")
		temporary_table_owner = input("temporary table owner =")
		data_schema_name = input("data schema name (schema containing observation_fact)=")
		option = input("totalnum information: counting observations or patients? (o/p) =")

		if metadata_schema_name != "" and metadata_table_name != "" and temporary_table_owner != "" and data_schema_name != "" and option != "":
			break
		print("All those fields are necessary")

#depending on the chosen option we count by number of patient related to a concept or by number of observation related to a concept.
subject = "patient" if option == "p" else "observation"
observation_attribute = "patient_num" if option == "p" else "text_search_index"

#replace terms in query
query = postgreSQL_script.format(
	metadata_schema_name=metadata_schema_name, metadata_table_name=metadata_table_name,
	temporary_table_owner=temporary_table_owner, data_schema_name=data_schema_name,
	subject=subject, observation_id=observation_attribute)


sqlScriptFilename = "updateTotalnum.psql"
with open(sqlScriptFilename, "w") as sqlFile:
	sqlFile.write(query)

print(f"Wrote generated postgreSQL script to {sqlScriptFilename}")