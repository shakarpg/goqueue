run:
	go run cmd/main.go

test:
	go test ./tests -v

build:
	go build -o goqueue cmd/main.go

docker-build:
	docker build -t goqueue:latest .

docker-run:
	docker run -p 8080:8080 goqueue:latest

clean:
	rm -f goqueue

.PHONY: run test build docker-build docker-run clean
