package aggregator

import (
	"fmt"
	"log"

	"github.com/net-reply-future-networks/k8s-log-aggregator/internal/utils"
)

type Aggregator struct {
	NatsClient utils.NatsClient
}

func (a *Aggregator) Run() {
	_, ch, err := a.NatsClient.GetSubscription()
	if err != nil {
		log.Fatal(err)
	}
	for {
		msg := <-ch
		fmt.Println(string(msg.Data))
	}
}
