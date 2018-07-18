FROM golang:alpine AS go-build-env


RUN apk add --no-cache git build-base

ADD . /go/src/github.com/lon9/discord-generalized-voice-bot/bot3
WORKDIR /go/src/github.com/lon9/discord-generalized-voice-bot/bot3

RUN go get
RUN go build -o /usr/bin/app


FROM alpine
RUN apk add --no-cache ca-certificates
COPY --from=go-build-env /usr/bin/app /app