 /* This script updates totalnum for all modifiers and concepts in the specified metadata table.
  *
  * The recursive formula for determining the totalnum value of a modifier is:
  * 	modifier totalnum = #(set of subjects directly referencing the modifier UNION of all set of subjects referencing the children of the modifier )
  *		
  * 	concept totalnum = #(set of subjects directly referencing the modifier UNION of all set of subjects referencing the children of the concept
  * 
  * A modifier is a parent of another modifier if the parent's c_fullname is the suffix of another child c_fullname.
  * A concept is a parent of a modifier if the parent concept c_fullname is equal to the m_applied_path column of the child.
  * A concept cannot have parents.
  * 
  * Hence one may have multiple layers of modifiers but there is a actually only one layer of concept. 
  * 
  * The script uses a bottom approach to computer this formula
  */

 /* TODO Create a script which takes as parameter the name of the database on which this script is going to be used and the tables that will be used. This script 
  * will then generate this postgreSQL script with the right values where needed
  */

--a temporary table used for computation of totalnum that links a modifer with a subject
CREATE TABLE i2b2metadata.concept_to_subject(
	patient_num integer NOT NULL,
	c_fullname character varying(2000) COLLATE pg_catalog."default",
	c_hlevel integer NOT NULL ,
	is_modifier boolean NOT NULL --if it is not a modifier it is a concept.
);


CREATE INDEX c_hlevel_cts
    ON i2b2metadata.concept_to_subject USING btree
    (c_hlevel ASC NULLS LAST)
    TABLESPACE pg_default;

CREATE INDEX c_fullname_idx_sphn_cts
    ON i2b2metadata.concept_to_subject USING btree
    (c_fullname COLLATE pg_catalog."default" ASC NULLS LAST)
    TABLESPACE pg_default;

ALTER TABLE i2b2metadata.concept_to_subject
    OWNER to i2b2metadata;



-- filling totalnum for leaf modifiers.
UPDATE i2b2metadata.sphn s
	SET c_totalnum = (  
		SELECT COUNT (DISTINCT o.patient_num ) FROM i2b2demodata.observation_fact o
		WHERE s.c_basecode = o.modifier_cd AND o.modifier_cd != '@'  -- '@' is for concepts
			AND o.encounter_num >= 0   -- counting only valid encounter

	)
WHERE s.c_facttablecolumn = 'MODIFIER_CD' AND s.c_basecode != '@';
				 

--In this code we insert into concept_to_subject the distinct (patient number, modfier fullname, modifier level) tuples that can be linked together by directly checking references between observation_fact and metadata
INSERT INTO i2b2metadata.concept_to_subject
SELECT DISTINCT o.patient_num, s.c_fullname, s.c_hlevel, TRUE
FROM i2b2demodata.observation_fact o, i2b2metadata.sphn s
WHERE o.encounter_num >= 0  AND s.c_basecode = o.modifier_cd 
	AND s.c_facttablecolumn = 'MODIFIER_CD' AND s.c_basecode != '@' 
	AND o.modifier_cd != '@';  -- '@' is for concepts 
 
--In this code we insert into concept_to_subject the distinct (patient number, concept fullname, concept level) tuples that can be linked together by directly checking references between observation_fact and metadata 
INSERT INTO i2b2metadata.concept_to_subject
SELECT DISTINCT o.patient_num, s.c_fullname, s.c_hlevel, FALSE
FROM i2b2demodata.observation_fact o, i2b2metadata.sphn s 
WHERE s.c_basecode = o.concept_cd AND o.modifier_cd = '@'  -- '@' is for concepts
	AND o.encounter_num >= 0   -- counting only valid encounter
	AND s.c_facttablecolumn = 'CONCEPT_CD';



/**
 * This code traverses iteratively the different heights of concepts/modfiers.
 * At each iteration the code adds the patients that relate to the child concept/modifiers to the set of patients that relate to the parent.
 * After having done that totalnum is updated for the parent of the current iteration. 
 */
CREATE OR REPLACE FUNCTION aggregateModifiersCounts() RETURNS void AS $$
	DECLARE
		min_depth int := (SELECT MIN(c_hlevel) FROM i2b2metadata.sphn s WHERE s.c_facttablecolumn = 'MODIFIER_CD');
		max_depth int := (SELECT MAX(c_hlevel) FROM i2b2metadata.sphn s WHERE s.c_facttablecolumn = 'MODIFIER_CD');
		before_min_depth int := min_depth + 1;
	BEGIN 
		-- the children height is height and the parent height is height-1.
		FOR height IN REVERSE max_depth..before_min_depth
		LOOP  

			--link the subjects of the child to the parents so those subject are linked with the parent when we count the distinct subjects linked to a parent.
			INSERT INTO i2b2metadata.concept_to_subject
			SELECT DISTINCT cts.patient_num, parent.c_fullname, (height-1), (parent.c_facttablecolumn = 'MODIFIER_CD')
			FROM i2b2metadata.sphn child, i2b2metadata.concept_to_subject cts, i2b2metadata.sphn parent
			WHERE child.c_totalnum > 0  AND cts.c_hlevel = height AND child.c_hlevel = height AND child.c_facttablecolumn = 'MODIFIER_CD' 
				AND cts.c_fullname = child.c_fullname AND
				(	
					(parent.c_hlevel = (height - 1) AND parent.c_facttablecolumn = 'MODIFIER_CD' AND --the parent is a modifier
						child.c_fullname LIKE (parent.c_fullname || '%') ESCAPE '|') -- Use | as an escape character instead of backslash which appears in windows paths.
					OR
					(child.m_applied_path = parent.c_fullname AND parent.c_facttablecolumn = 'CONCEPT_CD') --the parent is a concept
				)
			;
   			
			-- Then we count the number of distinct subjects that are linked to the parents and add that information to the totalnum column of the parent.
			UPDATE i2b2metadata.sphn s 
			SET c_totalnum = count_to_name.patient_count
			FROM ( 
				SELECT COUNT(DISTINCT cts.patient_num) AS patient_count, cts.c_fullname AS fullname
				FROM i2b2metadata.concept_to_subject cts
				WHERE  cts.is_modifier = TRUE AND cts.c_hlevel = (height-1)
				GROUP BY cts.c_fullname
			) AS count_to_name
			WHERE count_to_name.fullname = s.c_fullname;
			 
		END LOOP;
	END;
$$ LANGUAGE plpgsql;


SELECT aggregateModifiersCounts();

--update totalnum counts for concepts as we only did it for modifiers within the loop.
UPDATE i2b2metadata.sphn s 
SET c_totalnum = count_to_name.patient_count
FROM ( 
	SELECT COUNT(DISTINCT cts.patient_num) AS patient_count, cts.c_fullname AS fullname
	FROM i2b2metadata.concept_to_subject cts
	WHERE  cts.is_modifier = FALSE
	GROUP BY cts.c_fullname
) AS count_to_name
WHERE count_to_name.fullname = s.c_fullname;


--We are done with the temporary table.
DROP TABLE i2b2metadata.concept_to_subject;
 

 
