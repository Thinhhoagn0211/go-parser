run:
	docker compose -f ./docker-compose.yml --env-file=app.env up -d 

down:
	docker compose -f ./docker-compose.yml --env-file=app.env down

delete:
	docker compose -f ./docker-compose.yml --env-file=app.env down --remove-orphans -v

init:
	rm -r ./docs/
	swag init

build:
	docker build -t service-upload-media:v1.0.0 -f ./Dockerfile .
