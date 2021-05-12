# https://github.com/lcashdol/WAPot

#
# STAGE 1: build executable binary
#
FROM golang:1.16-buster as builder
WORKDIR /app

COPY . .
RUN CGO_ENABLED=0 go build

#
# STAGE 2: build a small image
#
FROM scratch
COPY --from=builder /app/wa-pot /bin/wa-pot
COPY html /var/www/html

ENTRYPOINT [ "/bin/wa-pot" ]
