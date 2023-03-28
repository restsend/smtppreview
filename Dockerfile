FROM golang:1.19-bullseye as builder
RUN mkdir /build 
ADD . /build/
WORKDIR /build 
ENV CGO_ENABLED=1 GO111MODULE=on GOPROXY=https://goproxy.cn
RUN go mod download
RUN go build -o smtppreview .

FROM ubuntu:22.04
ENV DEBIAN_FRONTEND noninteractive
ENV LANG C.UTF-8
RUN apt-get update && apt-get install -y ca-certificates tzdata
ENV HOME /app
WORKDIR /app

COPY --from=builder /build/smtppreview /app/
WORKDIR /app
CMD ["./smtppreview"]