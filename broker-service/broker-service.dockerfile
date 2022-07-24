#base of golang image
FROM golang:1.18-alpine as builder

RUN mkdir /app

COPY . /app

WORKDIR /app
#cgo is desabled i will not use it
RUN CGO_ENABLED=0 go build -o brokerApp ./cmd/api

#smaller image
FROM alpine:latest

RUN mkdir /app

COPY --from=builder /app/brokerApp /app

CMD [ "/app/brokerApp" ]