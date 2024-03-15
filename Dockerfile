############################
# STEP 0 get dependencies
############################
FROM golang:1.16.5 AS dependencies
ENV GOPRIVATE=*.countmax.ru
RUN apt-get update && apt-get install openssl -y
ENV DOMAIN_NAME=git.countmax.ru \
  TCP_PORT=443
RUN openssl s_client -connect $DOMAIN_NAME:$TCP_PORT -showcerts </dev/null 2>/dev/null | openssl x509 -outform PEM | tee /usr/local/share/ca-certificates/$DOMAIN_NAME.crt
RUN update-ca-certificates
WORKDIR /go/src
COPY go.mod .
COPY go.sum .
RUN go mod download
############################
# STEP 1 build executable binary
############################
FROM dependencies AS builder
LABEL maintainer="it@watcom.ru" version="0.0.5"
ARG BUILD_NUMBER
ARG GIT_HASH
ENV BUILD_NUMBER ${BUILD_NUMBER}
ENV GIT_HASH ${GIT_HASH}

ENV GOPRIVATE=*.countmax.ru
ENV GO111MODULE=on \
  CGO_ENABLED=0 \
  GOOS=linux \
  GOARCH=amd64

WORKDIR /go/src/
COPY . .

RUN go get -u github.com/swaggo/swag/cmd/swag
RUN swag init

RUN go build -ldflags="-X 'main.build=$BUILD_NUMBER' -X 'main.githash=$GIT_HASH'" -o /go/bin/layoutconfig.api
############################
# STEP 2 build a small image
############################
FROM alpine
RUN apk add --no-cache tzdata wget mailcap
ENV TZ=Europe/Moscow
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone
COPY --from=builder /go/bin/layoutconfig.api /go/bin/layoutconfig.api
COPY --from=builder /go/src/config.yaml /go/bin/config.yaml
COPY --from=builder /go/src/asset/index.html /go/bin/asset/index.html
WORKDIR /go/bin/
EXPOSE 8000 8001
CMD ["./layoutconfig.api"]
