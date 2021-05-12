#
# STAGE 1: build static web files
#
FROM node:14 as frontend
WORKDIR /src

#
# install dependencies
COPY client/package*.json ./
RUN npm install

#
# build client
COPY client/ .
RUN npm run build

#
# STAGE 2: build executable binary
#
FROM golang:1.16-buster as builder
WORKDIR /app

COPY server/. .
RUN go get -v -t -d .; \
    ./build

#
# STAGE 3: build a small image
#
FROM scratch
COPY --from=builder /app/bin/ioth /app/bin/ioth
COPY --from=frontend /src/dist/ /var/www

ENV DOCKER_API_VERSION=1.39
ENV IOTH_BIND=:8080
ENV IOTH_STATIC=/var/www

EXPOSE 8080

ENTRYPOINT [ "/app/bin/ioth" ]
CMD [ "serve" ]
