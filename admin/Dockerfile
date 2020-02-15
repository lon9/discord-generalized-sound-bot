FROM node:alpine as build-env

ARG BASE_URL=http://localhost:3000
ARG ADMIN_TITLE="Admin title here"

ENV NODE_ENV=production
ENV BASE_URL=${BASE_URL}
ENV ADMIN_TITLE=${ADMIN_TITLE}

WORKDIR /app
ADD . /app

RUN yarn install && yarn build

FROM node:alpine

ARG BASE_URL=http://localhost:3000
ARG ADMIN_TITLE="Admin title here"

ENV NODE_ENV=production
ENV BASE_URL=${BASE_URL}
ENV ADMIN_TITLE=${ADMIN_TITLE}

WORKDIR /app
COPY --from=build-env /app /app
