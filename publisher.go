package zros_go

import (
	"github.com/astaxie/beego/logs"
	"github.com/golang/protobuf/proto"
	"reflect"
)

type Publisher struct {
	topic       string
	messageType reflect.Type
	node        *defaultNode
	address     string
}

func NewPublisher(node *defaultNode, topic string, messageType reflect.Type) *Publisher {
	if messageType.Elem().Implements(msgType) {
		panic("NewPublisher messageType requires a proto")
	}
	publisher := &Publisher{}
	publisher.topic = topic
	publisher.node = node
	publisher.messageType = messageType
	return publisher
}

func (pub *Publisher) Publish(msg proto.Message) error {
	if reflect.TypeOf(msg) != pub.messageType {
		logs.Error("publish %s message failed message type must be %s", pub.topic, pub.messageType.String())
	}

	content, err := proto.Marshal(msg)
	if err != nil {
		panic("publish message is not a proto")
	}
	return pub.node.Publish(pub.topic, content)
}

func (pub *Publisher) SetAddress(address string) {
	pub.address = address
}

func (pub *Publisher) GetAddress() string {
	return pub.address
}

func (pub *Publisher) GetTopic() string {
	return pub.topic
}
