FROM golang AS builder

WORKDIR /home/app
COPY go.mod go.sum ./
RUN echo "download mod" \
    && go mod download

COPY . .
RUN echo "install app" \
    && make build

FROM golang
WORKDIR /home

# migration
COPY --from=builder /home/app/cmd/migrate/idiom.json ./idiom.json
COPY --from=builder /home/app/bin/migrate /usr/local/bin/

# zwei
COPY --from=builder /home/app/cmd/zwei/fonts ./fonts
COPY --from=builder /home/app/bin/zwei /usr/local/bin/

CMD [ "zwei" ]

