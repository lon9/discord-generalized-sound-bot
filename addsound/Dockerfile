FROM golang:alpine AS go-build-env

RUN apk add --no-cache git build-base
ADD . /go/src/github.com/lon9/discord-generalized-sound-bot/addsound
WORKDIR /go/src/github.com/lon9/discord-generalized-sound-bot/addsound

RUN go get
RUN go build -o /usr/bin/addsound

FROM alpine
RUN apk add --no-cache ffmpeg opus-dev liborvis-dev
COPY --from=go-build-env /usr/bin/addsound /usr/bin/addsound
