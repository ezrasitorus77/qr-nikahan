FROM golang:1.16-buster

ENV TZ Asia/Jakarta

RUN go get -d -v ./...
RUN go build -o /qr-nikahan

RUN rm -rf /.git

EXPOSE 8080

CMD ["qr-nikahan"]