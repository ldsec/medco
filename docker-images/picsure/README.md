This docker image sets up PIC-SURE with a PostgreSQL database, and allows to rely on Keycloak for the authentication by enabling OAuth2 client authentication.

PICSURE_2_TOKEN: token to talk server-to-PICSURE2 / intra-picsure requests (e.g. aggregate to query)
(should be valid to be used, but not used here)

in query to resource, the BEARER_TOKEN key in resource credentials will have the value configured in DB