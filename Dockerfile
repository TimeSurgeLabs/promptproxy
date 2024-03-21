FROM golang:1.22-alpine AS builder

LABEL org.opencontainers.image.source=github.com/chand1012/pocketbase

WORKDIR /pb

COPY . .

# build the base executable in examples/base
RUN go build -o pocketbase

FROM alpine:latest as runner

WORKDIR /pb

COPY --from=builder /pb/pocketbase /pb/pocketbase

EXPOSE 8090

ENTRYPOINT [ "/pb/pocketbase" ]

CMD [ "serve", "--http=0.0.0.0:8090", "--dir", "/pb_data"]
