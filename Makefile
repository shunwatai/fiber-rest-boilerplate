.PHONY: dev build clean \
	docker-dev docker-dev-build docker-dev-up docker-dev-down docker-dev-log \
	docker-prod docker-prod-build docker-prod-up docker-prod-down docker-prod-log \
	migrate-up migrate-down \
	tw-watch

dev:
	air -c .air.toml

build: 
	go build -race -o ./build/fiber-api .

# ref: https://unix.stackexchange.com/a/669683
clean: 
	find ./build/ -type f -executable -delete

docker-dev: docker-dev-build docker-dev-up
docker-prod: docker-prod-build docker-prod-up

docker-dev-build: 
	docker-compose -f compose-dev.yaml build --build-arg UID=$$(id -u) fiber-api-dev

docker-prod-build: 
	docker-compose -f compose-prod.yaml build --build-arg UID=$$(id -u) fiber-api-prod

docker-dev-up: 
	docker-compose -f compose-dev.yaml up -d

docker-prod-up: 
	docker-compose -f compose-prod.yaml up -d

docker-dev-down: 
	docker-compose -f compose-dev.yaml down

docker-prod-down: 
	docker-compose -f compose-prod.yaml down

docker-dev-log: 
	docker-compose -f compose-dev.yaml logs -f fiber-api-dev 

docker-prod-log: 
	docker-compose -f compose-prod.yaml logs -f fiber-api-prod

migrate-up: 
	go run main.go migrate-up $(filter-out $@,$(MAKECMDGOALS))

migrate-down: 
	go run main.go migrate-down $(filter-out $@,$(MAKECMDGOALS))

rbmq-worker: 
	go run main.go run-rbmq-worker

tw-watch:
	npx tailwindcss -i ./web/static/css/input.css -o ./web/static/css/output.css --watch

%:
	@:
