# Distributed TCP Chat Engine 🚀

A high-performance, scalable, and distributed TCP chat infrastructure built with **Go**, **Redis**, and **Nginx**.

Unlike standard single-server chat applications, this project implements a distributed architecture capable of handling horizontal scaling. It uses Nginx as a Layer 4 load balancer to route raw TCP connections to multiple Go server instances, and leverages Redis Pub/Sub as a central message broker to ensure users can communicate seamlessly regardless of which physical server they are connected to.

## 🏗 Architecture Overview

1. **The Entry Point (Nginx):** Acts as a Layer 4 TCP Load Balancer listening on port `8000`. It transparently proxies incoming user connections to one of the available Go backend servers.
2. **The Nodes (Go Servers):** Stateless TCP servers running on ports `9000`, `9001`, and `9002`. They handle raw socket connections, decode custom protocol packets, and manage the I/O for each user.
3. **The Brain (Redis Pub/Sub):** Acts as the global message bus. When User A sends a message to User B, the Go server publishes the payload to Redis. The server currently hosting User B listens to that Redis channel and flushes the message down User B's TCP socket.

## ✨ Key Features

* **Distributed State:** Servers are stateless; connection states are decoupled using Redis, allowing servers to crash or scale without dropping the entire network.
* **Custom TCP Protocol:** Implements a strict, pipe-delimited packet protocol (`TARGET_ID|SERVER_NO|FD|PAYLOAD`) to ensure data integrity over raw sockets.
* **Concurrency:** Highly concurrent design utilizing Go's goroutines for non-blocking I/O (independent Transmitter and Receiver loops per connection).
* **Layer 4 Load Balancing:** Uses Nginx `stream` modules to balance raw TCP streams rather than standard HTTP requests.

## 🛠 Prerequisites

To run this project locally, you need the following installed:
* [Go](https://go.dev/) (v1.24+)
* [Redis Server](https://redis.io/)
* [Nginx](https://nginx.org/)
* `nc` (Netcat) for testing client connections

## 🚀 How to Run

**1. Start the external services:**
Ensure Redis is running in the background:
```bash
sudo service redis-server start
