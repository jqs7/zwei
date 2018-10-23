FROM golang:1.11-alpine AS builder
RUN apk add --no-cache git make 
WORKDIR /home/app
COPY go.mod go.sum ./
RUN echo "download mod" \
    && go mod download

COPY . .
RUN echo "build app" \
    && make build-linux

FROM alpine
# RUN apk add --no-cache ca-
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /home/app/bin/zwei_unix_amd64 ./
ENTRYPOINT ["./zwei_unix_amd64"]
