package main

import (
	"L0/pkg/nats"
	"L0/pkg/server"
	"L0/pkg/storage"
	"github.com/jackc/pgx"
	"github.com/nats-io/stan.go"
	"log"
)

func main() {
	natsConn, err := stan.Connect("test-cluster", "1")
	if err != nil {
		log.Fatal("Failed to connect to NTA Streaming Server")
	}
	defer natsConn.Close()

	pgConn, err := pgx.Connect(pgx.ConnConfig{
		Host:     "127.0.0.1",
		Port:     5432,
		Database: "demo_service",
		User:     "postgres",
		Password: "password"})
	if err != nil {
		log.Fatal("Failed to connect to PostgreSQL server")
	}
	defer pgConn.Close()

	repo := storage.InitOrderRepo(pgConn)
	NATSClient := nats.InitClient(natsConn, repo)

	err = NATSClient.Start()
	if err != nil {
		log.Fatal("Failed to subscribe to order channel")
	}
	defer NATSClient.Close()

	httpServer := server.InitServer(repo)
	httpServer.Start()
}
