FROM alpine
ADD ca-certificates.crt /etc/ssl/certs/

COPY pod-heartbeat /

CMD /pod-heartbeat
