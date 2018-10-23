FROM golang AS builder

WORKDIR /home/app
COPY go.mod go.sum ./
RUN echo "download mod" \
    && go mod download

COPY . .
RUN echo "install app" \
    && make install

CMD ["zwei"]

# FROM golang:alpine
# RUN apk add --no-cache gcc musl-dev
# COPY --from=builder /home/app/cmd/zwei/fonts ./fonts
# COPY --from=builder /go/bin/zwei /usr/local/bin/

