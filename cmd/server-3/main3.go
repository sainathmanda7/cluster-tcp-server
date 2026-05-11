package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"go-tcp/protocol"
	"net"
	"strings"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()
var rdb *redis.Client

func initRedis() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
		PoolSize: 100,
	})

	// just make sure redis is up
	if err := rdb.Ping(ctx).Err(); err != nil {
		panic("redis not reachable: " + err.Error())
	}

	fmt.Println("Redis connected successfully")
}

func main() {
	initRedis()

	port := flag.String("port", "9002", "server port")
	flag.Parse()

	ln, err := net.Listen("tcp", ":"+*port)
	if err != nil {
		panic(err)
	}
	defer ln.Close()

	fmt.Println("listening on", *port)

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("accept err:", err)
			continue
		}

		go handle(conn)
	}
}

func handle(conn net.Conn) {
	defer conn.Close()

	userID := conn.RemoteAddr().String()

	// greet user
	conn.Write([]byte("WELCOME|" + userID + "\n"))
	fmt.Println("user joined:", userID)

	// listen to this user's channel
	sub := rdb.Subscribe(ctx, userID)
	defer sub.Close()

	// push redis messages to tcp
	go func() {
		for msg := range sub.Channel() {
			fmt.Fprintf(conn, "MSG|%s\n", msg.Payload)
		}
	}()

	reader := bufio.NewReader(conn)

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("user left:", userID)
			return
		}

		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		msg, err := protocol.Decoder(line)
		if err != nil {
			conn.Write([]byte("ERR|bad format\n"))
			fmt.Println("decode failed:", err)
			continue
		}

		// msg.IP = target user
		err = rdb.Publish(ctx, msg.IP, msg.Payload).Err()
		if err != nil {
			fmt.Println("publish err:", err)
		}
	}
}