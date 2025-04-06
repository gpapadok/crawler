build:
	docker build -t crawler-seed -f Dockerfile.seed .
	docker build -t crawler-node -f Dockerfile .

init:
	scripts/start_containers.sh

run:
	docker run --network crawler-net --rm crawler-node
