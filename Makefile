compile:
	go build -o ./build/ ./...

dev:
	go run ./...

exec:
	./build/queue

run: compile exec

	