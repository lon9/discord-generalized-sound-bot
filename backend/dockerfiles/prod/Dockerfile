FROM golang:alpine AS go-build-env

RUN apk add --no-cache git build-base

ADD . /go/src/github.com/lon9/discord-generalized-sound-bot/backend
WORKDIR /go/src/github.com/lon9/discord-generalized-sound-bot/backend

RUN go get
RUN go build -o /usr/bin/app

FROM alpine
RUN apk add --no-cache ffmpeg opus-dev libvorbis-dev
COPY --from=go-build-env /usr/bin/app /app
