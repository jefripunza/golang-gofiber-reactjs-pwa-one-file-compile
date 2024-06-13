# --- 🛠️ nodejs as builder --- #
FROM node:18-alpine as fe-builder
LABEL org.opencontainers.image.authors="jefriherditriyanto@gmail.com"

# make the folder
RUN mkdir /react-build
WORKDIR /react-build
COPY . .

# Install the dependencies
RUN yarn install

# Build the project and copy the files
RUN yarn build

# =============================================================================== #

# --- 🛠️ golang as builder --- #
FROM golang:1.22-alpine as be-builder

#-> Setup Environment
# ENV GOPATH /go
# ENV PATH $PATH:$GOPATH/bin
ENV GO111MODULE on
ENV CGO_ENABLED 0
ENV GOOS linux
ENV GOARCH amd64
ENV CGO 0

#-> 🌊 Install Require
RUN apk add --no-cache \
    gcc \
    musl-dev \
    tzdata

WORKDIR /build
RUN mkdir ./dist
COPY . .
COPY --from=fe-builder /react-build/dist/ ./dist

#-> 🌊 Install Golang Module
RUN go mod tidy

# 💯 Configuration
RUN sed -i 's/127.0.0.1:/:/g' /build/server/http/server.http.go
RUN sed -i 's/localhost/host.docker.internal/g' /build/server/env/mongodb.env.go
RUN sed -i 's/localhost/host.docker.internal/g' /build/server/env/rabbitmq.env.go

#-> ⚒️ Build App
RUN go build -o ./run

# =============================================================================== #

# --- 🚀 Finishing --- #
FROM alpine:latest as runner
WORKDIR /app

# Add the community repository to get ffmpeg
RUN echo "http://dl-cdn.alpinelinux.org/alpine/edge/community" >> /etc/apk/repositories

# Install ffmpeg along with the other tools
RUN apk add --no-cache openssl curl nano ffmpeg

COPY --from=be-builder /build/run     /app/run

# 💯 Last Configuration
# COPY --from=be-builder /build/.env    /app/.env
# RUN sed -i 's/localhost/host.docker.internal/g' .env

# RUN cp /usr/share/zoneinfo/Asia/Jakarta /etc/localtime \
#     && echo "Asia/Jakarta" > /etc/timezone

RUN chmod +x ./run

ENTRYPOINT ["/app/run"]
CMD ["run"]
