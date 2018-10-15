#!/usr/bin/env bash
set -Eeuo pipefail

cat > "$LIGHTTPD_WEB_ROOT/i2b2-admin/i2b2_config_data.js" <<EOL
{
    urlProxy: "index.php",
    urlFramework: "js-i2b2/",

    lstDomains: [ {
        domain: "$I2B2_DOMAIN_NAME",
        name: "Domain $I2B2_DOMAIN_NAME",
        urlCellPM: "http://i2b2:8080/i2b2/services/PMService/",
        allowAnalysis: true,
        adminOnly: true,
        debug: false
    } ]
}
EOL

cat > "$LIGHTTPD_WEB_ROOT/i2b2-client/i2b2_config_data.js" <<EOL
{
    urlProxy: "index.php",
    urlFramework: "js-i2b2/",

    lstDomains: [  {
        domain: "$I2B2_DOMAIN_NAME",
        name: "Domain $I2B2_DOMAIN_NAME",
        urlCellPM: "http://i2b2:8080/i2b2/services/PMService/",
        allowAnalysis: true,
        debug: false
    } ]
}
EOL

# webclients whitelist URLs
sed -i "s/\+\\\.\[a-zA-Z\]{2,5}//" "$LIGHTTPD_WEB_ROOT/i2b2-client/index.php"
sed -i "s/\+\\\.\[a-zA-Z\]{2,5}//" "$LIGHTTPD_WEB_ROOT/i2b2-admin/index.php"


# TODO: GENOMIC STUFF BACKEND
#sed -i "s#SHRINE_ONT_DB#$I2B2_MEDCO_DB_NAME#g" "$LIGHTTPD_WEB_ROOT/shrine-client/js-i2b2/cells/plugins/MedCo/php/sqlConnection.php"
#sed -i "s#SHRINE_ONT_USER#genomic_annotations#g" "$LIGHTTPD_WEB_ROOT/shrine-client/js-i2b2/cells/plugins/MedCo/php/sqlConnection.php"
#sed -i "s#SHRINE_ONT_PW#$DB_PASSWORD#g" "$LIGHTTPD_WEB_ROOT/shrine-client/js-i2b2/cells/plugins/MedCo/php/sqlConnection.php"
#
