package zros_go

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

func (node *defaultNode) AdvertiseService(service string, reqType Message, resType Message, callback interface{}) (ServiceServer, error) {
	return nil, nil
}

func (node *defaultNode) ServiceClient(service string, reqType Message, resType Message) (ServiceClient, error) {
	return nil, nil
}

func (node *defaultNode) Advertise(topic string, msgType Message) (Publisher, error) {
	return nil, nil
}

func (node *defaultNode) Subscriber(topic string, msgType Message, callback interface{}) (Subscriber, error) {
	return nil, nil
}