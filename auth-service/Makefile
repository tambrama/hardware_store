.PHONY: all run migrations

all:

run:
	go run ./cmd/auth

migrations:
	goose -dir ./migrations sqlite3 ./storage/auth.db up