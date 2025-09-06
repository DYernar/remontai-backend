FROM golang:1.23.1 AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN swag init -g main.go --parseDependency --parseInternal --exclude vendor,internal/config,internal/repository,internal/service

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o service main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=build /app/service .

COPY --from=build /app/internal/config/config.yaml .
COPY --from=build /app/docs /root/docs

EXPOSE 8080

CMD ["./service"]