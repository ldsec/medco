{
  urlProxy: "/shrine-proxy/request",
        urlFramework: "js-i2b2/",
        loginTimeout: 15, // in seconds
        username_label:"MedCo username:",
        password_label:"MedCo password:",
        lstDomains: [
                {
                    domain: "SHRINE_WEBCLIENT_DOMAIN",
                    name: "SHRINE_WEBCLIENT_NAME",
                    debug: true,
                    allowAnalysis: true,
                    urlCellPM: "http://i2b2-server:8080/i2b2/services/PMService/",
                    isSHRINE: true
                }
        ]
}
