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

# This is a script which is used to generate a postgresql script that will update the c_totalnum value of a concept/modifier
# so that c_totalnum contains the number of patients with observations (or the number of observations, depending on your choice) for the concept/modifier and its children.

#tip: The generate psql script is easier to read than the postgreSQL_script that is waiting to be formatted by the execution of this script.



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
  */

/**
 * This code traverses iteratively the different hierarchical levels of the modifiers.
 * At each iteration the code adds the {subject}s that relate to the modifiers to the set of {subject}s that relate to the parent (concept or modifier).
 * After having done that totalnum is updated for the parent of the current iteration.
 */
CREATE OR REPLACE FUNCTION aggregateModifiersCounts() RETURNS void AS $$
	DECLARE
		min_level int := (SELECT MIN(c_hlevel) FROM {metadata_schema_name}.{metadata_table_name} s WHERE LOWER(s.c_facttablecolumn) = 'modifier_cd');
		max_level int := (SELECT MAX(c_hlevel) FROM {metadata_schema_name}.{metadata_table_name} s WHERE LOWER(s.c_facttablecolumn) = 'modifier_cd');
	BEGIN

		raise info 'before for loop';
		-- the children level is level and the parent level is level-1.
		FOR level IN REVERSE max_level..min_level
		LOOP
			raise info 'we are at level %', level;

			--link the {subject}s of the modifier children to their parent modifers or concept.
			INSERT INTO {metadata_schema_name}.concept_to_id
				SELECT DISTINCT cti.identifier, parent.c_fullname, (level-1)
				FROM {metadata_schema_name}.{metadata_table_name} child, {metadata_schema_name}.concept_to_id cti, {metadata_schema_name}.{metadata_table_name} parent
				WHERE cti.c_hlevel = level AND child.c_hlevel = level AND LOWER(child.c_facttablecolumn) = 'modifier_cd'
					AND cti.c_fullname = child.c_fullname AND
					(
						( --the parent is a modifier at the level just above.
							parent.c_hlevel = (level - 1) AND LOWER(parent.c_facttablecolumn) = 'modifier_cd' AND
							child.c_fullname LIKE (parent.c_fullname || '%') ESCAPE '|'
							-- Use | as an escape character instead of the backslash which appears in windows paths and breaks sql queries.
						)
						OR
						( --the parent is a concept
							LOWER(parent.c_facttablecolumn) = 'concept_cd' AND parent.c_fullname LIKE child.m_applied_path ESCAPE '|'
						)
					)
			;


		END LOOP;
	END;
$$ LANGUAGE plpgsql;

/**
 * This code traverses iteratively the different levels of concepts.
 * At each iteration the code adds the {subject}s that relate to the concepts to the set of {subject}s that relate to the parent concept.
 * After having done that totalnum is updated for the parents of the current iteration.
 */
CREATE OR REPLACE FUNCTION aggregateConceptsCounts() RETURNS void AS $$
	DECLARE
		min_level int := (SELECT MIN(c_hlevel) FROM {metadata_schema_name}.{metadata_table_name} s WHERE LOWER(s.c_facttablecolumn) = 'concept_cd');
		max_level int := (SELECT MAX(c_hlevel) FROM {metadata_schema_name}.{metadata_table_name} s WHERE LOWER(s.c_facttablecolumn) = 'concept_cd');
	BEGIN

		raise info 'before for loop';
		-- the children level is level and the parent level is level-1.
		FOR level IN REVERSE max_level..min_level
		LOOP
			raise info 'we are at level %', level;

			--link the {subject}s of the modifier children to their parent modifers or concept.
			INSERT INTO {metadata_schema_name}.concept_to_id
				SELECT DISTINCT cti.identifier, parent.c_fullname, (level-1)
				FROM {metadata_schema_name}.{metadata_table_name} child, {metadata_schema_name}.concept_to_id cti, {metadata_schema_name}.{metadata_table_name} parent
				WHERE cti.c_hlevel = level AND child.c_hlevel = level AND LOWER(child.c_facttablecolumn) = 'concept_cd'
					AND cti.c_fullname = child.c_fullname
					--the parent is a concept at the level just above.
					AND parent.c_hlevel = (level - 1) AND LOWER(parent.c_facttablecolumn) = 'concept_cd' AND
					child.c_fullname LIKE (parent.c_fullname || '%') ESCAPE '|'
					-- Use | as an escape character instead of the backslash which appears in windows paths and breaks sql queries.
			;

		END LOOP;
	END;
$$ LANGUAGE plpgsql;

--create a temporary index on the lower case version of the c_facttablecolumn. This is useful to optimize the performances of the '= LOWER(c_facttablecolumn)' condition present in this script
CREATE INDEX index_factcolumn_lower ON {metadata_schema_name}.{metadata_table_name} ((lower(c_facttablecolumn)));

DO $$
BEGIN
raise info 'created temporary index on lower cased c_facttablecolumn column of {metadata_schema_name}.{metadata_table_name}';

DROP TABLE IF EXISTS {metadata_schema_name}.concept_to_id;

--a temporary table used for computation of totalnum that links a concept/modifier with a {subject}
CREATE TABLE {metadata_schema_name}.concept_to_id(
	identifier integer NOT NULL, -- the identifier of the {subject}
	c_fullname character varying(2000) COLLATE pg_catalog."default", -- the full name of the concept or modifier.
	c_hlevel integer NOT NULL -- the level of the concept/modifier
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
WHERE LOWER(s.c_facttablecolumn) = 'modifier_cd' AND s.c_basecode != '@';

raise info 'Done filling totalnum counting observations or patients directly linked to the modifiers';

--In this code we insert into concept_to_id the distinct (patient/observation ID, modifier fullname, modifier level) tuples that can be linked together by directly checking references between observation_fact and metadata
INSERT INTO {metadata_schema_name}.concept_to_id
SELECT DISTINCT o.{observation_id}, s.c_fullname, s.c_hlevel
FROM {data_schema_name}.observation_fact o, {metadata_schema_name}.{metadata_table_name} s
WHERE s.c_basecode = o.modifier_cd
	AND LOWER(s.c_facttablecolumn) = 'modifier_cd' AND s.c_basecode != '@'
	AND o.modifier_cd != '@';  -- '@' is for concepts

--In this code we insert into concept_to_id the distinct (patient/observation ID, concept fullname, concept level) tuples that can be linked together by directly checking references between observation_fact and metadata
INSERT INTO {metadata_schema_name}.concept_to_id
SELECT DISTINCT o.{observation_id}, s.c_fullname, s.c_hlevel
FROM {data_schema_name}.observation_fact o, {metadata_schema_name}.{metadata_table_name} s
WHERE s.c_basecode = o.concept_cd AND o.modifier_cd = '@'  -- '@' is for concepts
	AND LOWER(s.c_facttablecolumn) = 'concept_cd';


PERFORM aggregateModifiersCounts();

raise info 'After call to aggregate modifiers {subject} function.';

PERFORM aggregateConceptsCounts();

raise info 'After call to aggregate concepts {subject} function.';

UPDATE {metadata_schema_name}.{metadata_table_name} s
SET c_totalnum = count_to_name.count -- count is either the count of distinct patients or the count of distinct observations per concept, modifier depending on what you chose to count upon.
FROM (
	SELECT COUNT(DISTINCT cti.identifier) AS count, cti.c_fullname AS fullname
	FROM {metadata_schema_name}.concept_to_id cti
	GROUP BY cti.c_fullname
) AS count_to_name
WHERE count_to_name.fullname = s.c_fullname;


--We are done with the temporary table.
DROP TABLE {metadata_schema_name}.concept_to_id;


--no need to explicitly commit since this is a implicit pgsql anonymous function which automatically commit the results at the END statement.
END;
$$ LANGUAGE plpgsql;

DROP INDEX IF EXISTS index_factcolumn_lower;

"""

print("""
This script will generate a postgreSQL script which you will be able to execute on your database in order to update
the totalnum column of your metadata/ontology table. We need a few informations about your database in order to generate that
SQL script:
_ The ontolgy schema name, that is the name of the schema that contains your metadata table.
_ The name of your metadata table which holds information about the modifiers and concepts such as the totalnum column.
_ The owner of the temporary table that will be used in order to perform the count query.
_ The data schema name, that is, the name of the schema which contains your observation_fact table.

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
