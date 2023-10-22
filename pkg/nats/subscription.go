package nats

import (
	"L0/pkg/model"
	"L0/pkg/storage"
	"bytes"
	"encoding/json"
	"github.com/nats-io/stan.go"
	"log"
)

// Client - клиент NATS Streaming, отвечающий за подписку и получение данных с канала
type Client struct {
	conn          stan.Conn
	repo          storage.OrderRepo
	subscriptions []stan.Subscription
}

func InitClient(conn stan.Conn, repo storage.OrderRepo) *Client {
	return &Client{conn: conn, repo: repo}
}

// Start - метод, запускающий подписки на каналы и считывание данных из них
func (c *Client) Start() error {
	sub, err := c.conn.Subscribe("channel", func(m *stan.Msg) {
		dec := json.NewDecoder(bytes.NewReader(m.Data))
		dec.DisallowUnknownFields()
		var order model.Order
		if err := dec.Decode(&order); err != nil {
			log.Println("Error. Failed to unmarshall received data\n" + err.Error())
		} else {
			if err := c.repo.Insert(order); err != nil {
				log.Println("Error. Failed to insert new order.\n" + err.Error())
			}
		}
	})
	c.subscriptions = append(c.subscriptions, sub)
	return err
}

// Close - отписывание от всех каналов и закрытие соединения NATS
func (c *Client) Close() {
	for _, sub := range c.subscriptions {
		sub.Unsubscribe()
	}
	c.conn.Close()
}
