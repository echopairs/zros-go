package zros_go

import (
	"reflect"
	pb "zros-go/zros_rpc"

	"github.com/astaxie/beego/logs"
	"github.com/golang/protobuf/proto"
)

type ServiceClient struct {
	serviceName string
	reqType     reflect.Type
	resType     reflect.Type
	node        *defaultNode
	ready       bool
}

func NewServiceClient(node *defaultNode, service string, reqType reflect.Type, resType reflect.Type) *ServiceClient {
	if resType.Elem().Implements(msgType) || resType.Elem().Implements(msgType) {
		panic("NewServiceClient reqType and resType requires a proto")
	}
	client := &ServiceClient{}
	client.serviceName = service
	client.reqType = reqType
	client.resType = resType
	client.node = node
	client.ready = false
	return client
}

func (sc *ServiceClient) Call(request proto.Message, timeout int) (proto.Message, *pb.Status) {
	if reflect.TypeOf(request) != sc.reqType {
		logs.Error("call  %s service failed request type must be %s", sc.serviceName, sc.reqType.String())
		panic("request must be match")
	}

	status := new(pb.Status)
	if !sc.ready {
		status.Code = pb.Status_UNKNOWN
		status.Details = "server not ready"
		return nil, status
	}
	content, err := proto.Marshal(request)
	if err != nil {
		panic("request message is not a proto")
	}
	status.Code = pb.Status_OK
	response, err := sc.node.Call(sc.serviceName, content, timeout)
	if err != nil {
		logs.Warn("call rpc %s failed for %s", sc.serviceName, err.Error())
		status.Code = pb.Status_UNKNOWN
		status.Details = err.Error()
		return nil, status
	}

	value := reflect.New(sc.resType.Elem())
	out, ok := value.Interface().(proto.Message)
	if !ok {
		logs.Error("convert failed")
	}
	err = proto.Unmarshal(response.ResponseData, out)
	if err != nil {
		logs.Error("proto unmarshal failed")
		status.Code = pb.Status_UNKNOWN
		status.Details = err.Error()
		return nil, status
	}
	return out, status
}

func (sc *ServiceClient) GetServiceName() string {
	return sc.serviceName
}
