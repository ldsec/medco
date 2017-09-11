#!/bin/bash
set -e

# meant to be called by Dockerfile of i2b2-web
# env var used: I2B2_DOMAIN_NAME, LIGHTTPD_WEB_ROOT

cat > "$LIGHTTPD_WEB_ROOT/i2b2-admin/i2b2_config_data.js" <<EOL
{
    urlProxy: "index.php",
    urlFramework: "js-i2b2/",

    lstDomains: [ {
        domain: "$I2B2_DOMAIN_NAME",
        name: "Domain $I2B2_DOMAIN_NAME",
        urlCellPM: "http://i2b2-server:8080/i2b2/services/PMService/",
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
        urlCellPM: "http://i2b2-server:8080/i2b2/services/PMService/",
        allowAnalysis: true,
        debug: false
    } ]
}
EOL

cat > "$LIGHTTPD_WEB_ROOT/medco-i2b2-client/i2b2_config_data.js" <<EOL
{
    urlProxy: "index.php",
    urlFramework: "js-i2b2/",

    lstDomains: [  {
        domain: "$I2B2_DOMAIN_NAME",
        name: "Domain $I2B2_DOMAIN_NAME",
        urlCellPM: "http://i2b2-server:8080/i2b2/services/PMService/",
        allowAnalysis: true,
        debug: false
    } ]
}
EOL

cat > "$LIGHTTPD_WEB_ROOT/index.html" <<EOL
<html><head><title>I2b2-web</title></head><body>

<div align="center">

<p><a href="/i2b2-admin">I2b2 admin</a></p>
<p><a href="/i2b2-client">I2b2 client</a></p>
<p><a href="/phppgadmin">PhpPgAdmin</a></p>
<p><a href="http://localhost:9990">WildFly Management</a></p>
<p><a href="http://localhost:8080/i2b2">I2b2 Axis Management</a></p>
<p><a href="/medco-i2b2-client">MedCo-I2b2 client (temporary)</a></p>
</div>

</body>
</html>
EOL

# i2b2 client and admin whitelist URL
sed -i "s/\"http:\/\/localhost\"/\"http:\/\/i2b2-server:8080\"/" "$LIGHTTPD_WEB_ROOT/i2b2-admin/index.php"
sed -i "s/\"http:\/\/localhost\"/\"http:\/\/i2b2-server:8080\"/" "$LIGHTTPD_WEB_ROOT/i2b2-client/index.php"
# todo: tmp
sed -i "s/\"http:\/\/localhost\"/\"http:\/\/i2b2-server:8080\"/" "$LIGHTTPD_WEB_ROOT/medco-i2b2-client/index.php"

