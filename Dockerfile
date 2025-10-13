FROM golang:1.25-alpine AS builder

WORKDIR /app
RUN apk add make
COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go install github.com/google/wire/cmd/wire@latest
RUN make wire

# Install migrate with postgres support
RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# REMOVE THIS LINE - don't run migrations during build
# RUN make migrate-up

RUN go build -o cupping.backend.app ./cmd/app/

FROM alpine:latest AS application

WORKDIR /app

# Copy the built application
COPY --from=builder /app/cupping.backend.app ./

# Copy migrate binary from builder for runtime migrations
COPY --from=builder /go/bin/migrate /usr/local/bin/migrate

# copy Makefile for runtime migrate-up target
COPY --from=builder /app/Makefile ./Makefile

# include email templates used at runtime
COPY --from=builder /app/templates ./templates

# include migrations used at runtime
COPY --from=builder /app/migrations ./migrations

# On container start, run migrations via Makefile, then start app
CMD ["/bin/sh", "-c", "make migrate-up && ./cupping.backend.app"]