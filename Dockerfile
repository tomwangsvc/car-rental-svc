FROM scratch

COPY ca-certificates.crt /etc/ssl/certs/
COPY car-svc /
COPY schema/* /
COPY zoneinfo.zip /

CMD ["/car-svc"]
