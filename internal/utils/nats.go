package utils

import (
	"github.com/nats-io/nats.go"
)

type NatsClient struct {
	client *nats.Conn
}

func (nc *NatsClient) GetInstance() (*nats.Conn, error) {
	if nc.client == nil {
		client, err := nats.Connect("nats:4222")
		if err != nil {
			return nil, err
		}
		nc.client = client
	}
	return nc.client, nil
}

func (nc *NatsClient) Publish(data []byte) error {
	client, err := nc.GetInstance()
	if err != nil {
		return err
	}
	err = client.Publish("LOGS", data)

	return err
}

func (nc *NatsClient) GetSubscription() (*nats.Subscription, chan *nats.Msg, error) {
	client, err := nc.GetInstance()
	if err != nil {
		return nil, nil, err
	}
	ch := make(chan *nats.Msg, 64)
	sub, err := client.ChanSubscribe("LOGS", ch)
	return sub, ch, err
}
