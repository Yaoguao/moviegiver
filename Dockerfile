FROM golang:1.23-alpine AS builder

WORKDIR /app

RUN apk --no-cache add bash make gcc git musl-dev

# install dependency
COPY go.mod go.sum ./
RUN go mod download

# copy common code
COPY . .

RUN go build -o /app/bin ./cmd/api

RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest


# BUILD APP ------------------------------------------------------------------------------------------------

FROM alpine AS app

COPY --from=builder /app/bin/api /


CMD ["/api"]

# BUILD MIGRATOR ------------------------------------------------------------------------------------------

FROM alpine AS migrator

WORKDIR /app

RUN apk add --no-cache postgresql-client

COPY --from=builder /go/bin/migrate /usr/local/bin/migrate
COPY --from=builder /app/migrations /app/migrations

COPY /bin/sh/migrate.sh migrate.sh
RUN chmod +x /app/migrate.sh
ENTRYPOINT ["/app/migrate.sh"]

#ENTRYPOINT ["/bin/sh", "-c", "migrate -path=/app/migrations -database=$MOVIEGIVER_DB_DSN up"]


