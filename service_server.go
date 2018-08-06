package zros_go

import (
	pb "zros-go/zros_rpc"
	"reflect"
	"github.com/golang/protobuf/proto"
	"fmt"
)

type ServiceServer struct {
	serviceName  string
	serviceCb    reflect.Value
	reqType 	 reflect.Type
	resType 	 reflect.Type
	node 		 *defaultNode
}

func NewServiceServer(node *defaultNode, service string, reqType reflect.Type, resType reflect.Type, serviceCb interface{}) *ServiceServer{

	if resType.Elem().Implements(msgType) || resType.Elem().Implements(msgType) {
		panic("NesServiceServer reqType and resType requires a proto")
	}
	// check serviceCb
	fv := reflect.ValueOf(serviceCb)
	ft := fv.Type()
	if ft.Kind() != reflect.Func {
		panic(fmt.Sprintf("NewServiceServer serviceCb requires a function"))
	}
	server := &ServiceServer{}
	server.serviceName = service
	server.reqType = reqType
	server.resType = resType
	server.serviceCb = reflect.ValueOf(serviceCb)
	server.node = node
	return server
}

func (ss *ServiceServer) GetServiceName() string {
	return ss.serviceName
}

func (ss *ServiceServer) Invoke(request *pb.ServiceRequest) (*pb.ServiceResponse, error) {
	requestData := request.RequestData
	response := new(pb.ServiceResponse)
	var in []reflect.Value
	var iv reflect.Value
	if ss.reqType.Implements(msgType) {
		iv = reflect.New(ss.reqType.Elem())		// new value
		iv.Elem().Set(reflect.Zero(ss.reqType.Elem()))
		in = append(in, iv)
	} else {
		iv = reflect.New(ss.reqType)
		in = append(in, iv.Elem())
	}
	imsg := (iv.Interface()).(proto.Message)
	err := proto.Unmarshal(requestData, imsg)
	if err != nil {
		status := &pb.Status{}

		status.Code = pb.Status_INVALID_ARGUMENT
		status.Details = err.Error()
		response.Status = status
		return response, nil
	}
	responseData := ss.serviceCb.Call(in)
	var res interface{}
	if ss.resType.Implements(msgType) {
		ov := reflect.New(ss.resType)
		ov.Elem().Set(responseData[0])
		res = ov.Elem().Interface()
	} else {
		ov := reflect.New(ss.resType)
		ov.Elem().Set(responseData[0])
		res = ov.Interface()
	}
	msg, ok := res.(proto.Message)
	if !ok {
		panic("result is not a proto")
	}
	response.ResponseData, err = proto.Marshal(msg)
	if err != nil {
		response.Status.Code = pb.Status_UNKNOWN
		response.Status.Details = err.Error()
	}
	return response, nil
}