build-containers:
	docker build . -t cristovaoolegario/tasks-api -t latest
	docker tag cristovaoolegario/tasks-api ghcr.io/cristovaoolegario/tasks-api:latest
	docker build -f Dockerfile.consumer . -t cristovaoolegario/tasks-consumer -t latest
	docker tag cristovaoolegario/tasks-consumer ghcr.io/cristovaoolegario/tasks-consumer:latest

start:
	docker-compose up -d

deploy:
	kubectl apply -f manifests/mysql/config.yaml
	kubectl apply -f manifests/mysql/secret.yaml
	kubectl apply -f manifests/mysql/mysql.yaml
	kubectl apply -f manifests/task-api/secret.yaml
	kubectl apply -f manifests/task-api/task-api.yaml
	minikube service task-api-service

swagger:
	swag init -g cmd/api/rest/main.go
