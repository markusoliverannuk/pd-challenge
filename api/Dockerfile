# Build stage
FROM --platform=linux/amd64 golang:alpine AS build

RUN apk add --no-cache gcc musl-dev

WORKDIR /app

COPY . .

RUN CGO_ENABLED=1 GOARCH=amd64 go build -o challenge ./cmd/main.go

# Final stage
FROM --platform=linux/amd64 alpine:latest

WORKDIR /app

COPY --from=build /app/challenge /app/

COPY --from=build /app/ /app/

RUN chmod +x /app/challenge

EXPOSE 8050

CMD ["/app/challenge"]