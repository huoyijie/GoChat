FROM ubuntu:latest

LABEL org.opencontainers.image.authors="huoyijie"

COPY target/gochat-server /usr/bin/

EXPOSE 8888

ENTRYPOINT ["gochat-server"]