FROM alpine

WORKDIR /usr/src

RUN apk add --no-cache bash nmap-ncat;

COPY *.sh ./

ENTRYPOINT ["./entrypoint.sh"]
