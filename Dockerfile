FROM golang:1.25-alpine AS builder

WORKDIR /app
RUN apk add make
COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go install github.com/google/wire/cmd/wire@latest
RUN make wire

RUN go install github.com/golang-migrate/migrate/v4/cmd/migrate@latest
RUN make migrate-up

RUN go build -o cupping.backend.app ./cmd/app/

FROM alpine:latest AS application

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