package main

import (
	zros "zros-go"
	pb "zros-go/zros_rpc"
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
	node := zros.NewDefaultNode("example_service_client")

	// 3. create service client
	reqType := reflect.TypeOf(&zros_example.TestServiceRequest{})
	resType := reflect.TypeOf(&zros_example.TestServiceResponse{})
	client, _ := node.ServiceClient("test_service", reqType, resType)

	request := &zros_example.TestServiceRequest{}
	request.Detail = "hello world"
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for {
		select {
		case <- ticker.C:
			res, status := client.Call(request, -1)
			if status.Code != pb.Status_OK {
				logs.Info("rpc call failed, because of  %s", status.Details)
			} else {
				logs.Info("the res detail is %s", res.(*zros_example.TestServiceResponse).Detail)
			}
		}

	}
}