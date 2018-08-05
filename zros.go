package zros_go

import (
	"github.com/golang/protobuf/proto"
	"reflect"
)

var msgType = reflect.TypeOf((*proto.Message)(nil)).Elem()

var gsd *GrpcServiceDiscovery

func Init(masterAddress string) (err error){
	gsd, err = NewGrpcServiceDiscovery(masterAddress)
	if err != nil {
		return err
	}
	return nil
}

func GetGlobalServiceDiscovery () *GrpcServiceDiscovery {
	return gsd
}



type Publisher interface {
	Publisher(msg proto.Message)
}

type ServiceClient interface {
	Call() error
}

type IServiceServer interface {
	GetServiceName() string
}

type Subscriber interface {

}

type Node interface {
	Spin()
	AdvertiseService(service string, reqType proto.Message, resType proto.Message, callback interface{}) (ServiceServer, error)
	ServiceClient(service string, reqType proto.Message, resType proto.Message) (ServiceClient, error)
	Advertise(topic string, msgType proto.Message) (Publisher, error)
	Subscriber(topic string, msgType proto.Message, callback interface{}) (Subscriber, error)
}

