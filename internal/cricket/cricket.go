package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"
)

type Cricket struct {
	Interval int16
	Interupt bool
}

func (c *Cricket) Chirp() {
	fmt.Println("Starting chirp")
	for !c.Interupt {
		fmt.Println(os.Getpid(), c.GetRandomString(20))
		time.Sleep(time.Duration(c.Interval) * time.Second)
	}
}

func (c *Cricket) GetRandomString(length int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, length)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func main() {
	cricket := Cricket{
		Interval: 1,
	}
	cricket.Chirp()

	time.Sleep(time.Duration(10) * time.Second)
	cricket.Interupt = true
	time.Sleep(time.Duration(10) * time.Second)
}
