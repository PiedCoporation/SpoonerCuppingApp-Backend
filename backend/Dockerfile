FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN which wire || echo "wire not found"
RUN go install github.com/google/wire/cmd/wire@latest
RUN go generate ./...
RUN ls -la internal/infrastructures/wire/*/
RUN find internal/infrastructures/wire -name "*.go" -type f

RUN go build -o cupping.backend.app ./cmd/app/

FROM scratch

COPY --from=builder /app/cupping.backend.app .

CMD ["./cupping.backend.app"]