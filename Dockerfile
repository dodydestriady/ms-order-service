FROM golang:1.25.1 AS build

WORKDIR /app

RUN go install github.com/air-verse/air@latest

COPY go.mod go.sum ./

RUN go mod download

COPY . .

FROM build AS development

CMD ["air"]

FROM build AS production

RUN CGO_ENABLED=0 GOOS=linux go build -o main cmd/server/main.go

FROM alpine:latest
WORKDIR /root/

COPY --from=production /app/main .

EXPOSE 8080

CMD ["./main"]