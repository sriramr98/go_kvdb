build:
	go build -o bin/kvdb .

dev:
	leader &
	follower

leader:
	go run main.go --leader --port=8080 > logs/leader.log 2>&1

follower:
	go run main.go --port=8082 > logs/follower.log 2>&1