.PHONY: build run docker

build:
	go build -o cache-server ./cmd/cache-server

run: build
	./cache-server

docker-build:
	docker build -t cache-service:dev .

docker-run: docker-build
	docker run -p 8080:8080 cache-service:dev
