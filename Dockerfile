FROM --platform=linux/amd64 golang:1.21.0-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main.exe

FROM --platform=linux/amd64 alpine:latest

WORKDIR /app

COPY --from=builder /app/main.exe .

COPY .env ./

CMD ["./main.exe"]