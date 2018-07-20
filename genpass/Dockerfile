FROM golang:alpine AS go-build-env

RUN apk add --no-cache git

WORKDIR /go/src/github.com/lon9/discord-generalized-sound-bot/genpass
ADD . /go/src/github.com/lon9/discord-generalized-sound-bot/genpass
RUN go get
RUN go build -o /usr/bin/genpass

FROM alpine
COPY --from=go-build-env /usr/bin/genpass /usr/bin/genpass