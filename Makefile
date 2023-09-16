
all: sqlc tidy lint test build

build:
	go build -ldflags="-X 'github.com/MartyHub/mail-api/version.Version=development'" -race

clean:
	rm -fr .coverage.out mail-api db/gen/

start:
	podman-compose up -d >/dev/null 2>&1

stop:
	podman-compose down 2>/dev/null

lint:
	go vet ./...
	golangci-lint run

migrate:
	tern migrate --config db/tern_docker.conf --migrations db/migrations/

sqlc:
	sqlc generate --experimental --file db/sqlc.yaml
	goimports -w db/gen/querier.go

run:
	MAIL_API_DATABASE_URL=postgres://postgres:postgres@localhost:5432/postgres \
	MAIL_API_DEVELOPMENT=true \
	MAIL_API_SENDER_INTERVAL=10s \
	MAIL_API_SMTP_FROM=test@tests.org \
	MAIL_API_SMTP_PORT=1025 \
	./mail-api

test:
	gotest -coverprofile .coverage.out -race -timeout 10s

tidy:
	go mod tidy
