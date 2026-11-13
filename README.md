# Distributed TCP Chat Engine
<img width="1672" height="941" alt="High-Concurrency TCP Server Monitoring Dashboard" src="https://github.com/user-attachments/assets/6803df0d-31fb-484f-b808-168e35edd06f" />

A high-performance, scalable, and distributed TCP chat infrastructure built with **Go**, **Redis**, and **Nginx**.

Unlike standard single-server chat applications, this project implements a distributed architecture capable of handling horizontal scaling. It uses Nginx as a Layer 4 load balancer to route raw TCP connections to multiple Go server instances, and leverages Redis Pub/Sub as a central message broker to ensure users can communicate seamlessly regardless of which physical server they are connected to.

---

## Architecture Overview

1. **The Entry Point (Nginx):** Acts as a Layer 4 TCP Load Balancer listening on port `8000`. It transparently proxies incoming user connections to one of the available Go backend servers.
2. **The Nodes (Go Servers):** Stateless TCP servers running on ports `9000`, `9001`, and `9002`. They handle raw socket connections, decode custom protocol packets, and manage the I/O for each user.
3. **The Brain (Redis Pub/Sub):** Acts as the global message bus. When User A sends a message to User B, the Go server publishes the payload to Redis. The server currently hosting User B listens to that Redis channel and flushes the message down User B's TCP socket.

---

## Key Features

- **Distributed State:** Servers are stateless; connection states are decoupled using Redis, allowing servers to crash or scale without dropping the entire network.
- **Custom TCP Protocol:** Implements a strict, pipe-delimited packet protocol (`TARGET_ID|SERVER_NO|FD|PAYLOAD`) to ensure data integrity over raw sockets.
- **Concurrency:** Highly concurrent design utilizing Go's goroutines for non-blocking I/O (independent Transmitter and Receiver loops per connection).
- **Layer 4 Load Balancing:** Uses Nginx `stream` modules to balance raw TCP streams rather than standard HTTP requests.

---

## Prerequisites

To run this project locally, ensure you have the following installed:
- [Go](https://go.dev/) (v1.24+)
- [Redis Server](https://redis.io/)
- [Nginx](https://nginx.org/)
- `nc` (Netcat) for testing client connections

---

## Getting Started

### 1. Start Nginx (TCP Load Balancer)

Run Nginx using the custom configuration file provided in the project root:

```bash
sudo nginx -c $(pwd)/nginx.conf
```

2. Start the Go Server Cluster
Open separate terminal windows to run the server instances:

```bash
# Terminal 1
go run ./cmd/server-1/main1.go -port=9000

# Terminal 2
go run ./cmd/server-2/main2.go -port=9001

# Terminal 3
go run ./cmd/server-3/main3.go -port=9002
```

3. Connect Clients
Open new terminal windows for each user (e.g., Alice, Bob) and connect through the Nginx load balancer:
```bash
nc localhost 8000
```

Message Protocol
To send messages between users, use the custom TCP protocol format below:
```bash
TARGET_ID|SERVER_NO|FD|YOUR_MESSAGE
```

Example
If Bob’s ID is 127.0.0.1:48760, Alice can send him a message like this:
127.0.0.1:48760|1|0|Hello from across the cluster!

Project Structure
```bash
go-tcp/
├── cmd/
│   ├── server-1/
│   │   └── main1.go         # Node 1 entrypoint
│   ├── server-2/
│   │   └── main2.go         # Node 2 entrypoint
│   └── server-3/
│       └── main3.go         # Node 3 entrypoint
│
├── protocol/
│   └── protocol.go          # TCP packet encoder/decoder
│
├── nginx.conf               # Layer 4 TCP load balancer config
├── go.mod
└── go.sum
```
