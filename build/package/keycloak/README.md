# Keycloak Docker Image for MedCo

The image will load at first boot a working configuration for MedCo from the file `default-medco-realm.json`.
If this default working configuration needs to be regenerated, read the following sections.

## Full export of Keycloak configuration
This needs a running instance of Keycloak, more specifically with the profile `dev-local-3nodes`.

```bash
# run export (will bring up a second instance of keycloak in the container)
docker exec -it dev-local-3nodes_keycloak_1 keycloak/bin/standalone.sh \
-Djboss.socket.binding.port-offset=100 -Dkeycloak.migration.action=export \
-Dkeycloak.migration.provider=singleFile \
-Dkeycloak.migration.realmName=master \
-Dkeycloak.migration.usersExportStrategy=REALM_FILE \
-Dkeycloak.migration.file=/tmp/keycloak-realm.json

# hit <Ctrl-C> to stop the second instance just brought up

# copy the exported file
docker cp dev-local-3nodes_keycloak_1:/tmp/keycloak-realm.json ./default-medco-realm.json
```
