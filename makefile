.PHONY: up down create tidy run


run:
	go run main.go

up:
	docker-compose up -d

down: 
	docker-compose down

create:
	docker build \
		-f . \
		-t dankstats-api:v1 \
		.
		
tidy:
	go mod tidy
	go mod vendor

