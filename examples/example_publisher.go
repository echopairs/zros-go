package main

import (
	zros "zros-go"
	"github.com/astaxie/beego/logs"
	"reflect"
	"zros-go/zros_example"
	"time"
)

func main() {

	// 1. init & run service discovery
	err := zros.Init("localhost:23333")
	if err != nil {
		logs.Error("zros init failed for %v", err)
		panic(err)
	}

	// 2. create node
	node := zros.NewDefaultNode("example_publisher")

	// 3. create publisher
	messageType := reflect.TypeOf(&zros_example.TestMessage{})
	publisher, _ := node.Advertise("test_topic", messageType)

	message := &zros_example.TestMessage{}
	message.Detail = "hello world"
	message.Count = 7
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for {
		select {
		case <- ticker.C:
			publisher.Publish(message)
		}
	}
}