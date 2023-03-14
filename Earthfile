VERSION 0.6
FROM golang:1.19.3

build-web:
  FROM node:16-alpine3.17
  WORKDIR /build
  COPY web .
  RUN npm config set registry https://registry.npm.taobao.org
  RUN npm install
  RUN npm run build
  SAVE ARTIFACT dist

build-server:
  WORKDIR /build
  RUN apt update
  RUN apt install -y git
  RUN apt install -y binutils
  COPY +build-web/dist web/dist
  COPY go.mod .
  COPY *.go .
  COPY apiimpl apiimpl
  COPY config config
  COPY utils utils
  COPY config.yaml .
  COPY script script
  RUN go env -w GOPRIVATE="gitee.com,jihulab.com"
  RUN go env -w GOPROXY="https://goproxy.cn,direct"
  RUN go mod tidy
  RUN go build -o gateway
  RUN strip gateway
  SAVE ARTIFACT gateway

docker:
  ARG TAG="latest"
  WORKDIR /app
  COPY +build-server/gateway .
  COPY config.docker.yaml config.yaml
  COPY script/check.go script/check.go
  COPY example/coding_provider_codegeex.go script/coding_provider.go
  ENV GOPATH=/go
  ENV GOROOT=/go
  ENTRYPOINT /app/gateway run
  EXPOSE 8080
  SAVE IMAGE --push linksaas/ai-gateway:$TAG 
