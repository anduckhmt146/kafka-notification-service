dev: 
	go run main.go start --CONFIG_PATH=./config/local.yaml

docker-compose-up:
	sudo docker compose up -d

docker-compose-down:
	docker compose down -d

kafka-producer:
	docker exec -it 93fd8707c96e kafka-console-producer.sh --broker-list localhost:9092 --topic kafka