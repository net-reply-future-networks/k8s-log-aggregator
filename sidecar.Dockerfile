# syntax=docker/dockerfile:1

FROM golang:latest AS builder
RUN adduser -u 10001 user2 --disabled-password
WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .
RUN cd cmd/sidecar && CGO_ENABLED=0 go build -ldflags="-s -w" -o sidecar


FROM alpine:latest

COPY --from=builder /etc/passwd /etc/passwd
USER user2
COPY --from=builder app/cmd/sidecar/sidecar ./
EXPOSE 8080
ENTRYPOINT [ "./sidecar" ]