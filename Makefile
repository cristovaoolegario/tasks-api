build-container:
	docker build . -t cristovaoolegario/tasks-api -t latest

start:
	docker-compose up
