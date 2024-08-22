# 🗄️ Gredis

This project is a simplified implementation of an in-memory database inspired by Redis, built using Go. It was created following the tutorial series from [Build Redis from Scratch](https://www.build-redis-from-scratch.dev/en/introduction).

## ✨ Features

-   🖥️ Basic Redis-compatible server
-   🛠️ Supports SET, GET, HSET, HGET, and PING commands
-   📡 RESP (Redis Serialization Protocol) implementation
-   💾 Data persistence using AOF (Append-Only File)

## 📋 Prerequisites

Before you begin, ensure you have the following installed on your system:

-   🐹 Go (version 1.16 or later)
-   🔧 Redis CLI (for testing)

## 🚀 Setting Up the Project

1. Clone the repository:

    ```
    git clone https://github.com/evanmschultz/gredis.git
    cd gredis
    ```

2. Install Redis CLI (if not already installed):

    - On macOS:
        ```
        brew install redis
        ```
    - On Ubuntu:
        ```
        sudo apt-get install redis-tools
        ```

3. Ensure that the Redis server is not running, as our application will use the same port:
    - On macOS:
        ```
        brew services stop redis
        ```
    - On Ubuntu:
        ```
        sudo systemctl stop redis
        ```

## 🏃‍♂️ Running the Project

1. Build and run the server:

    ```
    go run *.go
    ```

2. The server will start and listen on port 6379.

3. In another terminal, use Redis CLI to connect to your server:

    ```
    redis-cli
    ```

4. You can now interact with your server using Redis commands:
    ```
    > SET name John
    OK
    > GET name
    "John"
    > PING
    PONG
    ```

## 📁 Project Structure

-   `main.go`: Contains the main server logic and connection handling.
-   `resp.go`: Implements the RESP protocol for serialization and deserialization.
-   `handler.go`: Contains the command handlers (SET, GET, HSET, HGET, PING).
-   `aof.go`: Implements the Append-Only File (AOF) for data persistence.

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## 📄 License

This project is open source and available under the [MIT License](LICENSE).

## 🙏 Acknowledgements

This project was inspired by and created following the tutorial series from [Build Redis from Scratch](https://www.build-redis-from-scratch.dev/en/introduction). Special thanks to the author for the comprehensive guide.
