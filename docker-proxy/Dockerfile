FROM golang:1.16-alpine as builder
WORKDIR /app

COPY go/ .
RUN ./build

FROM alpine

WORKDIR /usr/src

RUN apk add --no-cache bash grep nmap-ncat tcpdump supervisor; \
    mkdir -p /var/log/supervisor /var/log/capture;

COPY src/* ./

COPY --from=builder /app/docker-proxy /usr/bin/docker-proxy

ENTRYPOINT ["./entrypoint.sh"]
