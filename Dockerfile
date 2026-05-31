FROM golang:alpine AS builder
RUN mkdir /build
ADD . /build/
WORKDIR /build

ENV CGO_ENABLED=0

RUN go build -o gs-go-fiber .

FROM alpine:latest
RUN addgroup -g 1000 noroot
RUN adduser -u 1000 -G noroot -h /home/noroot -D noroot
RUN mkdir /home/noroot/app
WORKDIR /home/noroot/app
EXPOSE 4040
EXPOSE 80
COPY --from=builder /build/gs-go-fiber /home/noroot/app/
CMD ["./gs-go-fiber"]