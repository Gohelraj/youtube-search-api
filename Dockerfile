ARG NAME=youtube-search-api
ARG SOURCEROOT=/go/src/${NAME}

# Builder Image
FROM golang:1.18-alpine as builder

ARG NAME
ARG SOURCEROOT

COPY . ${SOURCEROOT}
WORKDIR ${SOURCEROOT}

RUN go mod vendor -v
RUN GOOS=linux go build -o bin/${NAME} cmd/main.go

# Runner Image
FROM alpine:latest
ARG NAME
ARG SOURCEROOT
WORKDIR /usr/bin

ENV WAIT_VERSION 2.7.2
ADD https://github.com/ufoscout/docker-compose-wait/releases/download/$WAIT_VERSION/wait /wait
RUN chmod +x /wait
RUN apk update && apk add bash && apk --no-cache add ca-certificates && rm -rf /var/cache/apk/*
COPY --from=builder ${SOURCEROOT}/bin/${NAME} /usr/bin

CMD /wait && youtube-search-api