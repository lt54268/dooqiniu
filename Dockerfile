FROM golang:1.23-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o dooqiniu-app .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/dooqiniu-app .

COPY --from=builder /app/.env . 

EXPOSE 9090

CMD ["./dooqiniu-app"]
