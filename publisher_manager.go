package zros_go

import (
	"sync"
	zmq "github.com/pebbe/zmq4"
	"github.com/astaxie/beego/logs"
	"errors"
)

type PublisherManager interface {
	RegisterPublisher(publisher *Publisher) (error)
	UnregisterPublisher(topic string) (error)
	Publish(topic string, content []byte) (error)
}

type PubStub interface {
	Publish([]byte) (error)
	GetAddress() string
	Stop()
}

type ZmqPubStub struct {
	address	string
	topic 	string
	sync.Mutex
	sock 	*zmq.Socket
}

func NewZmqPubStub(topic string) (*ZmqPubStub, string) {

	sock, err := zmq.NewSocket(zmq.PUB)
	if err != nil {
		panic("zmq NewSocket zmq.Pub err")
	}
	err = sock.Bind("tcp://*:*")
	if err != nil {
		panic("zmq.pub Bind err")
	}
	address, err := sock.GetLastEndpoint()
	if err != nil {
		panic("zmq.pub GetLastEndpoint err")
	}
	logs.Info("bind success on address %s ", address)

	stub := &ZmqPubStub{
		topic:topic,
		address:address,
		sock:sock,
	}
	return stub, address
}

func (stub *ZmqPubStub) GetAddress() string {
	return stub.address
}

func (stub *ZmqPubStub) Publish(content []byte) error {
	stub.Lock()
	defer stub.Unlock()
	_, err := stub.sock.SendBytes(content, 0)
	return err
}

func (stub *ZmqPubStub) Stop() {
	stub.sock.Unbind(stub.address)
}

type PublishersImpl struct {
	pubsMutex	sync.Mutex
	publishers 	map[string]PubStub
}

func NewPublishersImpl() *PublishersImpl {
	return &PublishersImpl{
		publishers:make(map[string]PubStub),
	}
}

func (impl *PublishersImpl)RegisterPublisher(publisher *Publisher) (error) {
	// 1. register to memory first to gen address
	impl.addToMemory(publisher)
	// 2. register to master
	gsd := GetGlobalServiceDiscovery()
	err := gsd.AddPublisher(publisher)
	if err != nil {
		impl.UnregisterPublisher(publisher.GetTopic())
		return err
	}
	logs.Info("register publisher %s to master success", publisher.GetTopic())
	return nil
}

func (impl *PublishersImpl) addToMemory(publisher *Publisher) {
	impl.pubsMutex.Lock()
	defer impl.pubsMutex.Unlock()
	topic := publisher.topic
	_, ok := impl.publishers[topic]
	if ok {
		logs.Error("%s publisher already register ", topic)
		panic("cannot register same publisher")
	}
	address := impl.createPublisher(topic)
	publisher.SetAddress(address)
}

func (impl *PublishersImpl) UnregisterPublisher(topic string) (error) {
	impl.pubsMutex.Lock()
	defer impl.pubsMutex.Unlock()
	stub, ok := impl.publishers[topic]
	if ok {
		stub.Stop()
		delete(impl.publishers, topic)
	}
	return nil
}

func (impl *PublishersImpl) Publish(topic string, content []byte) (error) {
	impl.pubsMutex.Lock()
	defer impl.pubsMutex.Unlock()
	pub, ok := impl.publishers[topic]
	if !ok {
		return errors.New("please check  if register " + topic + " publisher")
	}
	return pub.Publish(content)
}

func (impl *PublishersImpl) createPublisher(topic string) string {
	stub, address := NewZmqPubStub(topic)
	impl.publishers[topic] = stub
	return address
}