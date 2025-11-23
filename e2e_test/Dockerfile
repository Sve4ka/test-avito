FROM golang:1.24.2-alpine3.21

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

CMD ["go", "test", "-v", "-tags=e2e", "./e2e_test/tests/..."]
