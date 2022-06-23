package aggregator

import (
	"log"
	"sync"

	"github.com/net-reply-future-networks/k8s-log-aggregator/internal/utils"
)

type Aggregator struct {
	NatsClient utils.NatsClient
	LogWriter  LogWriter
}

func (a *Aggregator) Run() {
	wg := sync.WaitGroup{}
	wg.Add(1)
	go a.RunNatsService()
	wg.Wait()
}

func (a *Aggregator) RunNatsService() {
	_, ch, err := a.NatsClient.GetSubscription()
	if err != nil {
		log.Fatal(err)
	}
	if err = a.LogWriter.DbHandlers.Startup(); err != nil {
		log.Fatal(err)
	}
	for {
		msg := <-ch
		a.LogWriter.WriteLog(msg.Data)
	}
}
