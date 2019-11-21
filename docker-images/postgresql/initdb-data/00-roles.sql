CREATE ROLE i2b2 LOGIN PASSWORD 'i2b2';
CREATE ROLE keycloak LOGIN PASSWORD 'keycloak';
CREATE ROLE genomicannotations LOGIN PASSWORD 'genomicannotations';
ALTER USER i2b2 CREATEDB;
ALTER USER keycloak CREATEDB;
ALTER USER genomicannotations CREATEDB;
