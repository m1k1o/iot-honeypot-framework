FROM golang:1.16-buster-slim as server
WORKDIR /app

COPY . .

RUN go get -v -t -d .; \
    ./build

ENTRYPOINT [ "bin/ioth" ]
CMD [ "serve" ]
