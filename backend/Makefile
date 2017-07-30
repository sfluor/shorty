build:
	GOOS=linux go build -ldflags "-linkmode external -extldflags -static" -o shorty
	docker-compose build

dev:
	GOOS=linux go build -ldflags "-linkmode external -extldflags -static" -o shorty
	docker-compose build
	docker-compose up	