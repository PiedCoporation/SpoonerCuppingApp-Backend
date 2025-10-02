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

# Install curl for health checks
RUN apk --no-cache add curl

COPY --from=builder /app/cupping.backend.app .

CMD ["./cupping.backend.app"]