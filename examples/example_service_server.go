package main

import (
	zros "zros-go"
	"github.com/astaxie/beego/logs"
	"zros-go/zros_example"
	"reflect"
)

func ServiceCb(req *zros_example.TestServiceRequest) (*zros_example.TestServiceResponse) {
	return nil
}

func main() {

	// 1. init & run service discovery
	err := zros.Init("localhost:23333")
	if err != nil {
		logs.Error("zros init failed for %v", err)
	}

	// 2. create node
	node := zros.NewDefaultNode("example_service_server")

	// 3. create server
	reqType := reflect.TypeOf(zros_example.TestServiceRequest{})
	resType := reflect.TypeOf(zros_example.TestServiceResponse{})
	node.AdvertiseService("test_service", reqType, resType, ServiceCb)

	c := make(chan interface{})
	_ = <- c
}
