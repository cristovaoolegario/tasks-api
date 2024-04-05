build-containers:
	docker build . -t cristovaoolegario/tasks-api -t latest
	docker build -f Dockerfile.consumer . -t cristovaoolegario/tasks-consumer -t latest

start:
	docker-compose up -d
