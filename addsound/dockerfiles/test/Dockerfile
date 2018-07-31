FROM golang:alpine

ADD . /go/src/github.com/lon9/discord-generalized-sound-bot
WORKDIR /go/src/github.com/lon9/discord-generalized-sound-bot/addsound
RUN apk add --no-cache git build-base ffmpeg opus-dev libvorbis-dev
RUN go get -v -t && \
  rm -rf /go/src/github.com/lon9/discord-generalized-sound-bot