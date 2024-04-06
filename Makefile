build-containers:
	docker build . -t cristovaoolegario/tasks-api -t latest
	docker build -f Dockerfile.consumer . -t cristovaoolegario/tasks-consumer -t latest

start:
	docker-compose up -d

deploy:
	kubectl apply -f manifests/mysql/config.yaml
	kubectl apply -f manifests/mysql/secret.yaml
	kubectl apply -f manifests/mysql/mysql.yaml
	kubectl apply -f manifests/task-api/secret.yaml
	kubectl apply -f manifests/task-api/task-api.yaml
	minikube service task-api-service
