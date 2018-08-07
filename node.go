package zros_go

import (
	"reflect"
	"fmt"

	pb "zros-go/zros_rpc"
)

// *defaultNode implements Node interface
type defaultNode struct {
	NodeAddress 	string
	NodeName 		string
	ssm 			ServiceServerManager
	scm 			ServiceClientManager
}

func NewDefaultNode(nodeName string) *defaultNode {
	node := &defaultNode{
		NodeName:nodeName,
	}
	node.ssm = NewGrpcServerImpl()
	node.scm = NewServiceClientsImpl()
	node.Spin()
	return node
}

func (node *defaultNode) SetNodeAddress(nodeAddress string) {
	node.NodeAddress = nodeAddress
}

func (node *defaultNode) Spin() {
	realAddress, err := node.ssm.Start()
	if err != nil {
		panic(fmt.Sprintf("Node Spin failed for %s", err.Error()))
	}
	node.NodeAddress = realAddress
}

func (node *defaultNode) AdvertiseService(service string, reqType reflect.Type, resType reflect.Type, callback interface{}) (*ServiceServer, error) {
	if len(service) <= 0 {
		panic("AdvertiseService failed, service cannot be empty")
	}
	server := NewServiceServer(node, service, reqType, resType, callback)
	err := node.ssm.RegisterServer(server)
	if err != nil {
		return nil, err
	}
	return server, nil
}

func (node *defaultNode) ServiceClient(service string, reqType reflect.Type, resType reflect.Type) (*ServiceClient, error) {
	if len(service) <= 0 {
		panic("ServiceClient failed, service cannot be empty")
	}
	client := NewServiceClient(node, service, reqType, resType)
	err := node.scm.RegisterClient(client)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func (node *defaultNode) Advertise(topic string, msgType reflect.Type) (Publisher, error) {
	return nil, nil
}

func (node *defaultNode) Subscriber(topic string, msgType reflect.Type, callback interface{}) (Subscriber, error) {
	return nil, nil
}

func (node *defaultNode) Call(serviceName string , content []byte, timeout int) (*pb.ServiceResponse, error){
	return node.scm.Call(serviceName, content, timeout)
}