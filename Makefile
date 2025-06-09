build:
	@go build -o bin/shelter-it cmd/app/main.go

run: build
	@./bin/shelter-it

migration:
	@migrate create -ext sql -dir cmd/migrate/migrations $(filter-out $@, $(MAKECMDGOALS))

migrate-up:
	@go run cmd/migrate/main.go up

migrate-down:
	@go run cmd/migrate/main.go down

migrate-force:
	@go run cmd/migrate/main.go force

migrate-version:
	@go run cmd/migrate/main.go version
