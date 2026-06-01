FROM golang:alpine AS builder
RUN mkdir /build
ADD . /build/
WORKDIR /build

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

# Download dependencies
RUN go mod download

RUN go build -ldflags "-s -w" -o bank-api .

FROM alpine:latest

RUN addgroup -g 1000 noroot

RUN adduser -u 1000 -G noroot -h /home/noroot -D noroot

RUN mkdir /home/noroot/app

WORKDIR /home/noroot/app

ENV API_PORT=3540 \
  API_VERSION=v1 \
  DB_HOST=/tmp/bank.db \
  AWS_REGION=us-east-1

COPY --from=builder /build/bank-api /home/noroot/app/

CMD ["./bank-api"]