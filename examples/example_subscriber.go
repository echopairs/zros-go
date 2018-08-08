package main

import (
	zros "zros-go"
	"github.com/astaxie/beego/logs"
	"reflect"
	"zros-go/zros_example"
)

func SubscriberCb(message *zros_example.TestMessage) {
	logs.Info("receive message %s", message.Detail)
}

func main() {

	// 1. init & run service discovery
	err := zros.Init("localhost:23333")
	if err != nil {
		logs.Error("zros init failed for %v", err)
		panic(err)
	}

	// 2. create node
	node := zros.NewDefaultNode("example_subscriber")

	// 3. create subscriber
	messageType := reflect.TypeOf(&zros_example.TestMessage{})
	node.Subscriber("test_topic", messageType, SubscriberCb)

	c := make(chan interface{})
	_ = <- c
}