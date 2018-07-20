FROM golang:alpine

WORKDIR /go/src/github.com/lon9/discord-generalized-sound-bot/backend
RUN apk add --no-cache git build-base ffmpeg opus-dev libvorbis-dev
ADD . /go/src/github.com/lon9/discord-generalized-sound-bot/backend
RUN go get -v -t && \
  rm -rf /go/src/github.com/lon9/discord-generalized-sound-bot/backend