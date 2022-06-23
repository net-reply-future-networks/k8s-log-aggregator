package aggregator

import (
	"crypto/rand"
	"encoding/json"

	"github.com/net-reply-future-networks/k8s-log-aggregator/internal/sidecar"
)

type DbHandlers struct {
	Fh FileHandlers
	Hm Hashmap
}

type DbHandlersInterface interface {
	GetRecord(key string) ([]byte, error)
	SetRecord(key string, data []byte) error
	DeleteRecord(key string) error
	Startup() error
}

func (d *DbHandlers) GenerateKey() []byte {
	token := make([]byte, 16)
	rand.Read(token)
	return token
}

func (d *DbHandlers) GetRecord(key string) ([]byte, error) {
	start, ok := d.Hm.GetStart(key)
	if !ok {
		return nil, nil
	}
	return d.Fh.Retrieve(start)
}

func (d *DbHandlers) SetRecord(key string, data sidecar.Log) error {
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}
	start, err := d.Fh.Write(b)
	if err != nil {
		return err
	}
	hk := Hashkey{
		Key:   key,
		Start: start,
	}
	d.Hm = append(d.Hm, hk)
	return d.Fh.WriteHk(hk)
}

func (d *DbHandlers) DeleteRecord(key string) error {
	start, err := d.Fh.WriteDeletion()
	if err != nil {
		return err
	}
	hk := Hashkey{
		Key:   key,
		Start: start,
	}
	d.Hm = append(d.Hm, hk)
	return d.Fh.WriteHk(hk)
}

func (d *DbHandlers) Startup() error {
	Hm, err := d.Fh.ReadHm()
	if err != nil {
		return err
	}
	d.Hm = Hm
	return d.Fh.CreateDbIfNotExists()
}
