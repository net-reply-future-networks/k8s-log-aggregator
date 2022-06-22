# syntax=docker/dockerfile:1

FROM golang:latest AS builder
RUN adduser -u 10001 user --disabled-password
WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .
RUN cd internal/cricket && CGO_ENABLED=0 go build -ldflags="-s -w" -o cricket


FROM scratch

COPY --from=builder /etc/passwd /etc/passwd
USER user
COPY --from=builder app/internal/cricket/cricket ./
EXPOSE 8080
ENTRYPOINT [ "./cricket" ]