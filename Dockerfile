FROM golang:1.26-alpine AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o app .

FROM alpine:3.21

RUN apk add --no-cache ca-certificates

WORKDIR /app
COPY --from=builder /build/app ./app
COPY docker-entrypoint.sh ./

ENTRYPOINT ["./docker-entrypoint.sh"]
CMD ["./app"]