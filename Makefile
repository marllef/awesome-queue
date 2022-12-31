compile:
	go build -o ./build/ ./...

dev:
	go run ./...

exec:
	./build/consumer

run: compile exec

	