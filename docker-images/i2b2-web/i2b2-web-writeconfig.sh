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

cat > "$LIGHTTPD_WEB_ROOT/index.html" <<EOL
<html><head><title>I2b2-web</title>
<script>
document.addEventListener('click', function(event) {
  var target = event.target;
  if (target.tagName.toLowerCase() == 'a')
  {
      var port = target.getAttribute('href').match(/^:(\d+)(.*)/);
      if (port)
      {
         target.href = port[2];
         target.port = port[1];
      }
  }
}, false);
</script>
</head><body>

<div align="center">
<p><a href="/shrine-client">SHRINE client (MedCo)</a></p>
<p><br /><br /></p>

<p><a href="/i2b2-admin">I2b2 admin</a></p>
<p><a href="/i2b2-client">I2b2 client</a></p>
<p><a href="/phppgadmin">PhpPgAdmin</a></p>
<p><a href="/phpmyadmin">PhpMyAdmin</a></p>
<p><a href=":9990">WildFly Management</a></p>
<p><a href=":8080/i2b2">I2b2 Axis2 Management</a></p>
<p><a href=":6443/manager">Tomcat Management</a></p>
<p><a href=":6443/shrine-dashboard">SHRINE Dashboard</a></p>
<p><a href=":6443/steward">SHRINE Data Steward</a></p>
<p><a href="/shrine-webclient-update.php">Pull last MedCo Webclient commits</a></p>
</div>

</body>
</html>
EOL

cat > "$LIGHTTPD_WEB_ROOT/shrine-webclient-update.php" <<EOL
<?php
    echo '<html><head><title>Pull last commits?</title></head><body>';
    echo '<form><input type="submit" name="btnSubmit" value="Do it" /></form>';

    if (isset(\$_GET['btnSubmit']) or isset(\$_POST['btnSubmit'])) {
        \$message=shell_exec("/opt/shrine-webclient-update.sh 2>&1");
        echo '<p>';
        print_r(\$message);
        echo '</p>';
    }

    echo '</body></html>';
?>
EOL


cat > "$LIGHTTPD_WEB_ROOT/shrine-client/i2b2_config_data.js" <<EOL
{
  urlProxy: "index.php",
        urlFramework: "js-i2b2/",
        loginTimeout: 15, // in seconds
        username_label:"MedCo username:",
        password_label:"MedCo password:",
        lstDomains: [
                {
                    domain: "$I2B2_DOMAIN_NAME",
                    name: "Domain $I2B2_DOMAIN_NAME",
                    debug: true,
                    allowAnalysis: true,
                    urlCellPM: "http://i2b2-server:8080/i2b2/services/PMService/",
                    isSHRINE: true
                }
        ]
}
EOL

cat > "$LIGHTTPD_WEB_ROOT/shrine-client/js-i2b2/cells/SHRINE/cell_config_data.js" <<EOL
{
        files: [
                "SHRINE_ctrl.js",
                "i2b2_msgs.js"
        ],
        css: [],
        config: {
                name: "SHRINE Cell",
                description: "SHRINE Cell",
                category: ["core","cell","shrine"],
                newTopicURL: "/steward/client/index.html",
                readApprovedURL:"https://shrine-server:6443/shrine/rest/i2b2/request"
        }
}
EOL

cat > "/etc/lighttpd/conf-enabled/10-ssl.conf" <<EOL
\$SERVER["socket"] == "0.0.0.0:443" {
	ssl.engine  = "enable"
	ssl.ca-file = "$CONF_DIR/cacert.pem"
	ssl.pemfile = "$CONF_DIR/srv$NODE_IDX.pem"
    #todo: names in configuration profiles make more explicit

    # todo: enable + get ssl only
	# strict configuration from https://cipherli.st/
	#ssl.honor-cipher-order = "enable"
	#ssl.cipher-list = "EECDH+AESGCM:EDH+AESGCM:AES256+EECDH:AES256+EDH"
	#ssl.use-compression = "disable"
	#setenv.add-response-header = (
	#	"Strict-Transport-Security" => "max-age=15724800; includeSubdomains; preload",
	#	"X-Frame-Options" => "DENY",
	#	"X-Content-Type-Options" => "nosniff"
	#)

	#ssl.use-sslv2 = "disable"
	#ssl.use-sslv3 = "disable"

	# strict configuration from https://raymii.org/s/tutorials/Strong_SSL_Security_On_lighttpd.html
	#ssl.dh-file = "/etc/ssl/certs/dhparam.pem"
	#ssl.ec-curve = "secp384r1"
}
EOL

# webclients whitelist URLs
sed -i "s/\"http:\/\/localhost\"/\"http:\/\/i2b2-server:8080\"/" "$LIGHTTPD_WEB_ROOT/i2b2-admin/index.php"
sed -i "s/\"http:\/\/localhost\"/\"http:\/\/i2b2-server:8080\"/" "$LIGHTTPD_WEB_ROOT/i2b2-client/index.php"
sed -i "s/\"http:\/\/127.0.0.1\"/\"http:\/\/i2b2-server:8080\"/" "$LIGHTTPD_WEB_ROOT/shrine-client/index.php"
sed -i "s/\"http:\/\/localhost\"/\"https:\/\/shrine-server:6443\"/" "$LIGHTTPD_WEB_ROOT/shrine-client/index.php"

# shrine webclient fixes for integration in php environment
sed -i "s#default.htm#index.html#g" "$LIGHTTPD_WEB_ROOT/shrine-client/index.php"
sed -i '/CURLOPT_SSL_VERIFYPEER/i curl_setopt($proxyRequest, CURLOPT_SSL_VERIFYHOST, FALSE);' "$LIGHTTPD_WEB_ROOT/shrine-client/index.php"
sed -i "s#SHRINE_ONT_DB#$I2B2_DOMAIN_NAME#g" "$LIGHTTPD_WEB_ROOT/shrine-client/js-i2b2/cells/plugins/MedCo/php/sqlConnection.php"
sed -i "s#SHRINE_ONT_USER#shrine_ont#g" "$LIGHTTPD_WEB_ROOT/shrine-client/js-i2b2/cells/plugins/MedCo/php/sqlConnection.php"
sed -i "s#SHRINE_ONT_PW#$DB_PASSWORD#g" "$LIGHTTPD_WEB_ROOT/shrine-client/js-i2b2/cells/plugins/MedCo/php/sqlConnection.php"

