FROM golang:1.24.2-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./

COPY ./migrations ./migrations

RUN go mod download

COPY . .

RUN go build -o ./main ./cmd/main.go  && chmod +x ./main

FROM alpine:3.21


RUN apk add --no-cache ca-certificates tzdata

WORKDIR /root/
COPY --from=build /app/main .
COPY --from=build /app/migrations ./migrations

EXPOSE 8080
CMD ["./main"]
