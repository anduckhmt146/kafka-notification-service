# Kafka Notification Service

Implement Notification Service By Kafka Producer and Consumer.

P/s: Self-learning project for Kafka.

# Run Local

1. Run Kafka with Docker

```bash
    make docker-compose-up
```

Config:

- **Broker**: localhost:9092
- **Topic**: notification
- **Consumer Group**: kafka-notification-v1

2. Create MySQL Database

- Go to config/local.yaml and configure your account.

3. Run source code and migrate database

```bash
    make dev
```

4. Go to notification database, table users and create 2 mock users.

5. Test Kafka Producer & Consumer

- Produce Message

```bash
    POST: http://localhost:8080/message
```

```bash
{
    "fromID": 1,
    "toID": 2,
    "message": "Hello World"
}
```

- Consume Messsage

  ![Kafka Consumer](/assets/kafka_consumer.png)
