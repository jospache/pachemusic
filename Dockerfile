FROM golang:1.23-alpine as builder

WORKDIR /app

COPY . .

RUN go mod tidy
RUN go build -o /app/main ./cmd/main.go

FROM alpine:latest

WORKDIR /app

RUN apk add --no-cache \
    yt-dlp \
    ca-certificates \
    ffmpeg \
    bash

RUN mkdir /app/downloads

COPY --from=builder /app/main /app/main
COPY --from=builder /app/site /app/site

RUN touch /app/cookies.txt && chmod 644 /app/cookies.txt

#ENV API_KEY_YOUTUBE=your_api_key
#ENV CORS_ORIGIN=https://seu-frontend.example

EXPOSE 8080

CMD ["/app/main"]
