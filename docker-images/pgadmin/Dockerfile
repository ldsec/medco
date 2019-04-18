FROM dpage/pgadmin4:4.3

ENV SCRIPT_NAME="/pgadmin"
COPY reverseproxyfix.sh /
RUN chmod +x /reverseproxyfix.sh && \
    /reverseproxyfix.sh
