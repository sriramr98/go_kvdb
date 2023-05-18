build:
	go build -o bin/kvdb .

run: build
	./bin/kvdb

dev:
	go run main.go