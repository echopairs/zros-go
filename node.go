package zros_go

import (
	"reflect"
)

// *defaultNode implements Node interface
type defaultNode struct {
	NodeAddress 	string
	NodeName 		string
}

func NewDefaultNode(nodeName string) *defaultNode {
	return &defaultNode{
		NodeName:nodeName,
	}
}

func (node *defaultNode) SetNodeAddress(nodeAddress string) {
	node.NodeAddress = nodeAddress
}

func (node *defaultNode) Spin() {

}

func (node *defaultNode) AdvertiseService(service string, reqType reflect.Type, resType reflect.Type, callback interface{}) (*ServiceServer, error) {
	return nil, nil
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