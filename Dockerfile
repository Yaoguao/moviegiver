FROM golang:1.23

WORKDIR /app

# install dependency
COPY go.mod go.sum ./
RUN go mod download

# copy common code
COPY . .

# install go-migrate
RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# По идее не нужна данная строка
# COPY migrations /migrations

# build
RUN go build -o /app/bin ./cmd/api

RUN chmod +x /app/bin/api

CMD ["sh", "-c","migrate -path=./migrations -database=${GREENLIGHT_DB_DSN} up && /app/bin/api"]
