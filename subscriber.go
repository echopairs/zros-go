package zros_go

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/golang/protobuf/proto"
	"reflect"
)

type Subscriber struct {
	topic       string
	messageCb   reflect.Value
	messageType reflect.Type
	node        *defaultNode
}

func NewSubscriber(node *defaultNode, topic string, messageType reflect.Type, messageCb interface{}) *Subscriber {
	if messageType.Elem().Implements(msgType) {
		panic("NewSubscriber messageType requires a proto")
	}

	// check messageCb
	fv := reflect.ValueOf(messageCb)
	ft := fv.Type()
	if ft.Kind() != reflect.Func {
		panic(fmt.Sprintf("NewSubscriber messageCb requires a function"))
	}
	subscriber := &Subscriber{}
	subscriber.topic = topic
	subscriber.messageType = messageType
	subscriber.node = node
	subscriber.messageCb = reflect.ValueOf(messageCb)
	return subscriber
}

func (sub *Subscriber) HandleRawMessage(content []byte) {

	var in []reflect.Value
	var iv reflect.Value
	if sub.messageType.Implements(msgType) {
		iv = reflect.New(sub.messageType.Elem())
		iv.Elem().Set(reflect.Zero(sub.messageType.Elem()))
		in = append(in, iv)
	} else {
		// todo
	}
	imsg := (iv.Interface()).(proto.Message)
	err := proto.Unmarshal(content, imsg)
	if err != nil {
		logs.Warn("receive message Unmarshal err %s", err.Error())
		return
	}
	sub.messageCb.Call(in)
}

func (sub *Subscriber) GetTopic() string {
	return sub.topic
}
