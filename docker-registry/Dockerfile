FROM alpine

WORKDIR /home/user

COPY certs/domain.crt certs/domain.crt
COPY apply.sh apply.sh

ENTRYPOINT [ "./apply.sh" ]
