FROM golang:1.19-bullseye as builder
RUN mkdir /build 
ADD . /build/
WORKDIR /build 
RUN CGO_ENABLED=1 GO111MODULE=on go build -o smtppreview .

FROM ubuntu:22.04
ENV DEBIAN_FRONTEND noninteractive
ENV LANG C.UTF-8
RUN apt-get update && apt-get install -y ca-certificates tzdata
ENV HOME /app
WORKDIR /app

COPY --from=builder /build/smtppreview /app/
WORKDIR /app
CMD ["./smtppreview"]