FROM golang:1.25-alpine AS builder

WORKDIR /app
RUN apk add make
COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go install github.com/google/wire/cmd/wire@latest
RUN make wire

RUN go build -o cupping.backend.app ./cmd/app/

FROM alpine:latest AS application

WORKDIR /app

# Copy the built application
COPY --from=builder /app/cupping.backend.app ./

# include email templates used at runtime
COPY --from=builder /app/templates ./templates

# include migrations used at runtime
COPY --from=builder /app/migrations ./migrations

# Start app; migrations run in-app on startup
CMD ["./cupping.backend.app"]