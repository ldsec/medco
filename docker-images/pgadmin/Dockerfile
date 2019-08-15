FROM dpage/pgadmin4:4.11

ENV SCRIPT_NAME="/pgadmin"
COPY reverseproxyfix.sh /
RUN chmod +x /reverseproxyfix.sh && \
    /reverseproxyfix.sh
