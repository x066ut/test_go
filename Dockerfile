FROM golang:1.10
RUN mkdir -p /go/src/app/bin
WORKDIR /go/src/app/bin
EXPOSE 8082
ADD http_server /go/src/app/bin
ADD server.conf /go/src/app/bin
CMD ["chmod 644 /usr/bin/server.conf"]
CMD ["/go/src/app/bin/http_server"]