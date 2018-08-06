package zros_go

import (
	"reflect"
	"fmt"
)

// *defaultNode implements Node interface
type defaultNode struct {
	NodeAddress 	string
	NodeName 		string
	ssm 			ServiceServerManager
}

func NewDefaultNode(nodeName string) *defaultNode {
	node := &defaultNode{
		NodeName:nodeName,
	}
	node.ssm = NewGrpcServerImpl()
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
	return nil, nil
}

func (node *defaultNode) Advertise(topic string, msgType reflect.Type) (Publisher, error) {
	return nil, nil
}

func (node *defaultNode) Subscriber(topic string, msgType reflect.Type, callback interface{}) (Subscriber, error) {
	return nil, nil
}