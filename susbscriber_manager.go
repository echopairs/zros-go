package zros_go

import (
	zmq "github.com/pebbe/zmq4"
	"sync"

	"context"
	"github.com/astaxie/beego/logs"
	"reflect"
	pb "zros-go/zros_rpc"
)

type SubscriberManager interface {
	RegisterSubscriber(subscriber *Subscriber) error
	UnregisterSubscriber(topic string) error
}

type SubStub interface {
	Start()
	Stop()
	AddPublisher(address string) error
	RemovePublisher(address string) error
}

type ZmqSubStub struct {
	topic      string
	pubMutex   sync.Mutex
	pubAddress map[string]int
	sock       *zmq.Socket
	ctx        context.Context
	cancel     context.CancelFunc

	runMutex sync.Mutex
	running  bool

	callback reflect.Value
}

func NewZmqSubStub(topic string, callback reflect.Value) *ZmqSubStub {
	ctx, cancel := context.WithCancel(context.Background())
	sock, err := zmq.NewSocket(zmq.SUB)
	if err != nil {
		panic("New Socket zmq.SUB err")
	}

	return &ZmqSubStub{
		topic:      topic,
		callback:   callback,
		pubAddress: make(map[string]int),
		ctx:        ctx,
		cancel:     cancel,
		sock:       sock,
	}
}

func (stub *ZmqSubStub) Stop() {
	if stub.cancel != nil {
		stub.cancel()
		stub.running = false
	}
}

func (stub *ZmqSubStub) Start() {
	stub.runMutex.Lock()
	defer stub.runMutex.Unlock()
	if stub.running {
		logs.Info("the topic: %s was already receive message ", stub.topic)
		return
	}

	stub.running = true
	go stub.receiveMessage(stub.ctx)
}

func (stub *ZmqSubStub) AddPublisher(address string) error {
	stub.pubMutex.Lock()
	defer stub.pubMutex.Unlock()
	_, ok := stub.pubAddress[address]
	if !ok {
		stub.sock.Connect(address)
		stub.sock.SetSubscribe("")
		stub.pubAddress[address] = 1
	}
	stub.Start()
	return nil
}

func (stub *ZmqSubStub) RemovePublisher(address string) error {
	stub.pubMutex.Lock()
	defer stub.pubMutex.Unlock()
	_, ok := stub.pubAddress[address]
	if ok {
		stub.sock.Unbind(address)
		delete(stub.pubAddress, address)
	}
	return nil
}

func (stub *ZmqSubStub) receiveMessage(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			logs.Info("receive %s message quit", stub.topic)
			stub.sock.Close()
			return
		default:
			data, _ := stub.sock.RecvBytes(0)
			in := make([]reflect.Value, 1)
			in[0] = reflect.ValueOf(data)
			go stub.callback.Call(in)
		}

	}
}

type SubPair struct {
	sub  *Subscriber
	stub SubStub
}

type SubscribersImpl struct {
	subscriberMutex sync.Mutex
	subscribers     map[string]*SubPair
}

func NewSubscribersImpl() *SubscribersImpl {
	impl := &SubscribersImpl{
		subscribers: make(map[string]*SubPair),
	}
	gsd := GetGlobalServiceDiscovery()
	registerPublisherCallback := func(info *pb.PublisherInfo, status *pb.Status) error {
		return impl.RegisterPublisher(info, status)
	}
	gsd.SetRegisterPublisherCb(registerPublisherCallback)
	return impl
}

func (impl *SubscribersImpl) RegisterSubscriber(subscriber *Subscriber) error {

	// 1. register to master
	gsd := GetGlobalServiceDiscovery()
	err := gsd.AddSubscriber(subscriber)
	if err != nil {
		logs.Error("AddSubscriber to master failed for %s", err.Error())
		return err
	}

	// 2. register to memory when 1 ok
	impl.addToMemory(subscriber)
	logs.Info("register subscriber %s to master success", subscriber.GetTopic())
	return nil
}

func (impl *SubscribersImpl) addToMemory(subscriber *Subscriber) {
	impl.subscriberMutex.Lock()
	defer impl.subscriberMutex.Unlock()
	topic := subscriber.GetTopic()
	_, ok := impl.subscribers[topic]
	if ok {
		logs.Error("%s subscriber already register", topic)
		panic("cannot register same subscriber")
	}

	callback := func(in []byte) {
		subscriber.HandleRawMessage(in)
	}

	stub := NewZmqSubStub(topic, reflect.ValueOf(callback))
	impl.subscribers[topic] = &SubPair{
		subscriber,
		stub,
	}
}

func (impl *SubscribersImpl) UnregisterSubscriber(topic string) error {
	impl.subscriberMutex.Lock()
	defer impl.subscriberMutex.Unlock()
	sub, ok := impl.subscribers[topic]
	if ok {
		sub.stub.Stop()
		delete(impl.subscribers, topic)
	}
	return nil
}

func (impl *SubscribersImpl) RegisterPublisher(info *pb.PublisherInfo, status *pb.Status) error {
	impl.subscriberMutex.Lock()
	defer impl.subscriberMutex.Unlock()
	topic := info.Topic
	_, ok := impl.subscribers[topic]
	if !ok {
		status.Code = pb.Status_NOT_FOUND
		status.Details = topic + " publisher not register on"
		return nil
	}
	publishAddress := info.PhysicalNodeInfo.RealAddress
	pair := impl.subscribers[topic]
	pair.stub.AddPublisher(publishAddress)
	status.Code = pb.Status_OK
	return nil
}

func (impl *SubscribersImpl) UnregisterPublisher(info *pb.PublisherInfo, status *pb.Status) error {
	// todo
	return nil
}
