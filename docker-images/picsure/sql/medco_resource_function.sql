CREATE FUNCTION add_or_update_medco_resource(
  p_uuid uuid,
  p_description varchar,
  p_name varchar,
  p_resourcerspath varchar,
  p_targeturl varchar,
  p_token varchar
) RETURNS VOID AS $$
BEGIN
  -- update or insert based on name column

  UPDATE picsure_resource SET description=p_description, resourcerspath=p_resourcerspath,
                              targeturl=p_targeturl, token=p_token
    WHERE picsure_resource.name=p_name;

  INSERT INTO picsure_resource(uuid, description, name, resourcerspath, targeturl, token)
    SELECT p_uuid, p_description, p_name, p_resourcerspath, p_targeturl, p_token
    WHERE NOT EXISTS (SELECT 1 FROM picsure_resource WHERE picsure_resource.name=p_name);

END;
$$ LANGUAGE plpgsql;
