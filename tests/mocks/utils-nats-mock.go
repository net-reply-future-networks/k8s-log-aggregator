package mocks

import (
	"fmt"
	"testing"

	"github.com/nats-io/nats.go"
)

type NatsClientMock struct {
	MockPublish         MockPublish
	MockGetSubscription MockGetSubscription

	T *testing.T
}

type MockPublish struct {
	InData              []byte
	OutError            error
	ExpectedInvocations int
	ActualInvocations   int
}

type MockGetSubscription struct {
	OutSubscription *nats.Subscription
	OutChannel      chan *nats.Msg
	OutError        error
}

func (n *NatsClientMock) GetInstance() (*nats.Conn, error) {
	return nil, nil
}
func (n *NatsClientMock) Publish(data []byte) error {
	n.MockPublish.ActualInvocations++
	if string(data) != string(n.MockPublish.InData) {
		n.T.Error(fmt.Sprintf("(Publish) unexpected input data value, expected %s, got %s", string(n.MockPublish.InData), string(data)))
	}
	return n.MockPublish.OutError
}
func (n *NatsClientMock) GetSubscription() (*nats.Subscription, chan *nats.Msg, error) {
	return n.MockGetSubscription.OutSubscription, n.MockGetSubscription.OutChannel, n.MockGetSubscription.OutError
}
