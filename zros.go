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
	gsd.Spin()
	return gsd.IsConnectedToMaster()
}

func GetGlobalServiceDiscovery () *GrpcServiceDiscovery {
	return gsd
}

type Node interface {
	Spin()
	AdvertiseService(service string, reqType proto.Message, resType proto.Message, callback interface{}) (ServiceServer, error)
	ServiceClient(service string, reqType proto.Message, resType proto.Message) (ServiceClient, error)
	Advertise(topic string, msgType proto.Message) (Publisher, error)
	Subscriber(topic string, msgType proto.Message, callback interface{}) (Subscriber, error)
}

