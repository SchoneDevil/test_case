FROM golang:1.18 AS builder
WORKDIR /app
COPY go.* ./
RUN go mod download
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server .

FROM alpine:latest

COPY --from=builder /app/server /app/server
COPY --from=builder /app/config.yml /app/server

EXPOSE 7777
CMD  ["/app/server"]