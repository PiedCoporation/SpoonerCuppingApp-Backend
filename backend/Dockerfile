FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go generate ./...

RUN go build -o cupping.backend.app ./cmd/app/

FROM scratch

COPY --from=builder /app/cupping.backend.app .

CMD ["./cupping.backend.app"]