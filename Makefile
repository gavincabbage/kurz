integration: docker-dev
	docker-compose -f ./test/docker-compose.yml up

docker-dev:
	docker build -t kurz:dev -f test/Dockerfile .