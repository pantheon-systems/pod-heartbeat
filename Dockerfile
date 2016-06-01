FROM alpine

ADD https://curl.haxx.se/ca/cacert.pem /etc/ssl/certs/

COPY pod-heartbeat /

CMD /pod-heartbeat
