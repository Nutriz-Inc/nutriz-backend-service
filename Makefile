run:
	nodemon --watch . --ext go --exec "go run ./main.go"

test:
	ENV=test go test ./...

migrations:
	go run cmd/migrations/run.go cmd/migrations/seed.go