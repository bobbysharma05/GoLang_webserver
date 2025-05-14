# Kafka-to-Redis Integration with Web Interface

This Go application demonstrates a data pipeline that consumes messages from a Kafka topic, stores them in Redis, and provides HTTP endpoints to produce messages to Kafka (`POST /produce`) and retrieve data from Redis (`GET /data/{key}`). It showcases integration with Kafka, Redis, and a web server, built for simplicity and reliability using Docker.

## Project Overview

- **Purpose**: Process JSON messages (e.g., `{"id":"123","value":"Hello, Kafka!"}`) from Kafka, store them in Redis with their key (e.g., `123`), and allow users to send new messages or retrieve stored data via HTTP.
- **Technologies**:
  - **Go**: Application logic and HTTP server.
  - **Kafka**: Messaging system (topic: `test-topic`).
  - **Redis**: Key-value storage.
  - **Docker**: Runs Kafka and Redis containers.
- **Features**:
  - Consumes messages from Kafka and stores them in Redis.
  - HTTP `POST /produce` endpoint to send messages to Kafka.
  - HTTP `GET /data/{key}` endpoint to retrieve Redis data.
  - Logs message processing and storage for debugging.

## Prerequisites

- **Docker** and **Docker Compose**: To run Kafka and Redis containers.
- **Go**: Version 1.18 or higher for the application.
- **MinGW-w64**: Provides `gcc` for `cgo` compilation (required for `librdkafka`).
- **vcpkg**: Supplies `librdkafka` library for Kafka integration.
- **Windows**: Commands are tailored for Windows Command Prompt (`cmd`).

## Project Structure

```
webserver/
├── handlers/
│   └── server.go       # HTTP handlers for POST /produce and GET /data/{key}
├── kafka/
│   ├── config.go       # Kafka configuration (broker, topic)
│   ├── consumer.go     # Kafka consumer logic
│   └── producer.go     # Kafka producer logic
├── redis/
│   └── redis.go        # Redis client and storage functions
├── docker-compose.yml   # Defines Kafka and Redis services
├── go.mod              # Go module dependencies
├── go.sum              # Dependency checksums
├── main.go             # Application entry point
└── README.md           # This documentation
```

## Setup Instructions

Follow these steps to set up and run the project on Windows.

### 1. Clone or Prepare the Project

- Place the project in `C:\Users\sharm.LENOVO\go\src\webserver`.
- Alternatively, clone from a repository (if hosted):
  ```bash
  git clone <your-repo-url>
  cd webserver
  ```

### 2. Start Kafka and Redis

- Ensure Docker is running.
- Start containers:
  ```bash
  cd C:\Users\sharm.LENOVO\go\src\webserver
  docker-compose up -d
  ```
- Verify containers:
  ```bash
  docker ps
  ```
  Expected:
  ```
  CONTAINER ID   IMAGE                 COMMAND                  STATUS         PORTS                    NAMES
  <id>           apache/kafka:latest   "/__cacert_entrypoin…"   Up             0.0.0.0:9092->9092/tcp   kafka-broker
  <id>           redis:latest          "docker-entrypoint.s…"   Up             0.0.0.0:6379->6379/tcp   redis-test-instance
  ```

### 3. Create Kafka Topic

- Create the `test-topic` if it doesn’t exist:
  ```bash
  docker exec -it -w /opt/kafka/bin kafka-broker sh
  ./kafka-topics.sh --create --topic test-topic --bootstrap-server localhost:9092 --partitions 1 --replication-factor 1
  exit
  ```
- Verify:
  ```bash
  docker exec -it -w /opt/kafka/bin kafka-broker sh
  ./kafka-topics.sh --list --bootstrap-server localhost:9092
  exit
  ```

### 4. Configure Go Environment

- Set up `cgo` for Kafka library compilation:
  ```bash
  set CGO_ENABLED=1
  set CC=C:\msys64\mingw64\bin\gcc.exe
  set PATH=%PATH%;C:\msys64\mingw64\bin;C:\vcpkg
  ```
- Ensure MinGW-w64 and vcpkg are installed:
  - MinGW-w64: Run `C:\msys64\msys2_shell.cmd` and install `pacman -S mingw-w64-x86_64-gcc`.
  - vcpkg: Install `librdkafka` (e.g., `vcpkg install librdkafka`).

### 5. Run the Application

- Start the Go application:
  ```bash
  cd C:\Users\sharm.LENOVO\go\src\webserver
  go run main.go
  ```
- Expected output:
  ```
  Starting application...
  Starting HTTP server on :8080...
  Starting Kafka consumer...
  ```

## Testing the Application

Test the HTTP endpoints and verify data storage.

### 1. Send a Message (POST /produce)

- Use `curl` (available in MSYS2 or Windows 10+) or Postman:
  ```bash
  curl -X POST -H "Content-Type: application/json" -d "{\"id\":\"789\",\"value\":\"Test via HTTP\"}" http://localhost:8080/produce
  ```
- Expected response:
  ```json
  {"status":"Message produced"}
  ```
- Check application logs:
  ```
  Produced message to Kafka: {ID:789 Value:Test via HTTP}
  Stored data in Redis with key 789: {ID:789 Value:Test via HTTP}
  Retrieved data from Redis: {"id":"789","value":"Test via HTTP"}
  ```

### 2. Retrieve Data (GET /data/{key})

- Fetch data from Redis:
  ```bash
  curl http://localhost:8080/data/789
  ```
- Expected response:
  ```json
  {"data":"{\"id\":\"789\",\"value\":\"Test via HTTP\"}"}
  ```
- Test with another key:
  ```bash
  curl http://localhost:8080/data/123
  ```
  Expected (if data exists):
  ```json
  {"data":"{\"id\":\"123\",\"value\":\"Hello, Kafka!\"}"}
  ```
- If the key doesn’t exist:
  ```
  Key not found
  ```

### 3. Verify Redis Storage

- Check Redis directly:
  ```bash
  docker exec -it redis-test-instance redis-cli
  GET 789
  GET 123
  QUIT
  ```
- Expected:
  ```
  "{\"id\":\"789\",\"value\":\"Test via HTTP\"}"
  "{\"id\":\"123\",\"value\":\"Hello, Kafka!\"}"
  ```

### 4. Optional: Produce Messages via CLI

- Send test messages using Kafka’s console producer:
  ```bash
  docker exec -it -w /opt/kafka/bin kafka-broker sh
  echo '{"id":"123","value":"Hello, Kafka!"}' | ./kafka-console-producer.sh --bootstrap-server localhost:9092 --topic test-topic
  echo '{"id":"456","value":"Another message"}' | ./kafka-console-producer.sh --bootstrap-server localhost:9092 --topic test-topic
  exit
  ```
- Check logs for storage confirmation.

## Troubleshooting

- **Application Fails to Start**:
  - Verify Docker containers:
    ```bash
    docker ps
    docker logs kafka-broker
    docker logs redis-test-instance
    ```
  - Check `cgo` setup:
    ```bash
    gcc --version
    echo %CC%
    echo %CGO_ENABLED%
    ```
  - Run `go mod tidy` to sync dependencies.

- **HTTP Endpoints Fail**:
  - Ensure port `8080` is free:
    ```bash
    netstat -a -n -o | find "8080"
    ```
  - Test server:
    ```bash
    curl http://localhost:8080/produce
    ```
    Expected: `Method not allowed`.

- **Kafka/Redis Errors**:
  - Verify Kafka topic:
    ```bash
    docker exec -it -w /opt/kafka/bin kafka-broker sh
    ./kafka-topics.sh --list --bootstrap-server localhost:9092
    ```
  - Test Redis:
    ```bash
    docker exec -it redis-test-instance redis-cli PING
    ```

## Notes

- **Configurations**: Hardcoded for simplicity (`localhost:9092` for Kafka, `localhost:6379` for Redis, `:8080` for HTTP).
- **Module Path**: Assumes `webserver` (e.g., `import "webserver/kafka"`). If `github.com/sharm/webserver`, update imports accordingly.
- **Enhancements**:
  - Use environment variables for configurations.
  - Add graceful shutdown for the consumer and server.
  - Implement additional endpoints (e.g., `DELETE /data/{key}`).
