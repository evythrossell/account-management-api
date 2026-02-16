FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /account-management-api main.go

FROM scratch

COPY --from=builder /account-management-api /account-management-api

CMD ["/account-management-api"]