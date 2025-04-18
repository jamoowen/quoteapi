.PHONY: test build start-server clean


TEST_PKGS = ./internal/utils ./internal/cache/

bin:
	mkdir -p bin

# Create db directory if it doesn't exist
db:
	mkdir -p db

test:
	go test -v $(TEST_PKGS)

run:
	go run cmd/server/main.go

start-server:
	./bin/quoteapi

build:
	go build -o bin/quoteapi ./cmd/server/main.go

clean:
	rm -rf bin 

.PHONY: init-db
init-db: db
	sqlite3 db/quotedb.sqlite < migrations/001-setup.sql

.PHONY: recreate-dev-db
recreate-dev-db: db
	sqlite3 db/quotedb.sqlite < migrations/recreate-dev.sql


