FROM node:14-buster-slim as client
WORKDIR /app

#
# install dependencies
COPY package*.json ./
RUN npm install

#
# build client
COPY . .
RUN npm run build
