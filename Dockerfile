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

# Install curl for health checks, make for Makefile targets, and migrate binary
RUN apk --no-cache add curl bash ca-certificates libc6-compat make

# Install golang-migrate (download tarball and extract binary)
ENV MIGRATE_VERSION=v4.17.0
RUN curl -sSL -o /tmp/migrate.tgz \
      https://github.com/golang-migrate/migrate/releases/download/${MIGRATE_VERSION}/migrate.linux-amd64.tar.gz && \
    tar -xzf /tmp/migrate.tgz -C /usr/local/bin && \
    mv /usr/local/bin/migrate.linux-amd64 /usr/local/bin/migrate && \
    chmod +x /usr/local/bin/migrate && \
    rm -f /tmp/migrate.tgz

WORKDIR /app

COPY --from=builder /app/cupping.backend.app ./
# copy Makefile for runtime migrate-up target
COPY --from=builder /app/Makefile ./Makefile
# include email templates used at runtime
COPY --from=builder /app/templates ./templates
# include migrations used at runtime
COPY --from=builder /app/migrations ./migrations

# On container start, run migrations via Makefile, then start app
CMD ["/bin/sh", "-c", "make migrate-up && ./cupping.backend.app"]