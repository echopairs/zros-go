package zros_go

import (
	"context"
	"errors"
	"fmt"
	"github.com/astaxie/beego/logs"
	"google.golang.org/grpc"
	"sync"
	"time"
	pb "zros-go/zros_rpc"
)

var rpcStubCallShortTimeOut = 5 * 1000 * 1000
var rpcStubCallLongTimeOut = 3000 * 1000 * 1000

type ServiceClientManager interface {
	RegisterClient(client *ServiceClient) error
	UnregisterClient(serviceName string) error
	Call(serviceName string, content []byte, timeout int) (*pb.ServiceResponse, error)
}

type RpcStub interface {
	Call(serviceName string, content []byte, timeout int) (*pb.ServiceResponse, error)
}

type GrpcStub struct {
	nodeRpcStub pb.ServiceRPCClient
	conn        *grpc.ClientConn
	address     string
}

func NewGrpcStub(address string) (*GrpcStub, error) {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	c := pb.NewServiceRPCClient(conn)
	stub := &GrpcStub{
		nodeRpcStub: c,
		conn:        conn,
		address:     address,
	}
	return stub, nil
}

func (gs *GrpcStub) Call(serviceName string, content []byte, timeout int) (*pb.ServiceResponse, error) {
	if timeout == -1 {
		timeout = rpcStubCallLongTimeOut
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout))
	defer cancel()
	in := &pb.ServiceRequest{}
	in.ServiceName = serviceName
	in.RequestData = content
	return gs.nodeRpcStub.InvokeService(ctx, in)
}

type ServiceClientsImpl struct {
	clientsMutex sync.Mutex
	clients      map[string]*ServiceClient

	servicesMutex   sync.Mutex
	servicesAddress map[string]string

	stubsMutex sync.Mutex
	stubs      map[string]RpcStub
}

func NewServiceClientsImpl() *ServiceClientsImpl {
	impl := &ServiceClientsImpl{
		clients:         make(map[string]*ServiceClient),
		servicesAddress: make(map[string]string),
		stubs:           make(map[string]RpcStub),
	}
	gsd := GetGlobalServiceDiscovery()
	registerServiceServerCallback := func(info *pb.ServiceServerInfo, status *pb.Status) error {
		return impl.RegisterServiceClient(info, status)
	}
	gsd.SetRegisterServiceServerCb(registerServiceServerCallback)
	return impl
}

func (sci *ServiceClientsImpl) RegisterClient(client *ServiceClient) error {
	// 1. register to memory first
	sci.addToMemory(client)
	// 2. register to master
	gsd := GetGlobalServiceDiscovery()
	err := gsd.AddServiceClient(client)
	if err != nil {
		sci.UnregisterClient(client.GetServiceName())
		return err
	}
	logs.Info("register client %s to master success", client.GetServiceName())
	return nil
}

func (sci *ServiceClientsImpl) addToMemory(client *ServiceClient) {
	sci.clientsMutex.Lock()
	defer sci.clientsMutex.Unlock()
	serviceName := client.GetServiceName()
	_, ok := sci.clients[serviceName]
	if ok {
		logs.Warn("%s already register ", serviceName)
		panic("cannot register same client")
	}
	sci.clients[serviceName] = client
}

func (sci *ServiceClientsImpl) UnregisterClient(serviceName string) error {
	sci.clientsMutex.Lock()
	defer sci.clientsMutex.Unlock()
	_, ok := sci.clients[serviceName]
	if ok {
		delete(sci.clients, serviceName)
	}
	return nil
}

func (sci *ServiceClientsImpl) Call(serviceName string, content []byte, timeout int) (*pb.ServiceResponse, error) {
	sci.clientsMutex.Lock()
	defer sci.clientsMutex.Unlock()
	_, ok := sci.clients[serviceName]
	if !ok {
		return nil, errors.New("please check  if register this client")
	}
	stub, err := sci.GetRpcStub(serviceName)
	if err != nil {
		return nil, err
	}
	return stub.Call(serviceName, content, timeout)
}

func (sci *ServiceClientsImpl) RegisterServiceClient(serverInfo *pb.ServiceServerInfo, status *pb.Status) error {
	serviceName := serverInfo.ServiceName
	sci.clientsMutex.Lock()
	defer sci.clientsMutex.Unlock()
	_, ok := sci.clients[serviceName]
	if !ok {
		status.Code = pb.Status_NOT_FOUND
		status.Details = serviceName + "client not register on"
		return errors.New(status.Details)
	}
	sci.clients[serviceName].ready = true

	sci.servicesMutex.Lock()
	defer sci.servicesMutex.Unlock()
	sci.servicesAddress[serviceName] = serverInfo.PhysicalNodeInfo.RealAddress
	status.Code = pb.Status_OK
	return nil
}

func (sci *ServiceClientsImpl) UnregisterServiceClient(serverInfo *pb.ServiceServerInfo, status *pb.Status) error {
	// todo
	return nil
}

func (sci *ServiceClientsImpl) GetRpcStub(serviceName string) (RpcStub, error) {
	sci.servicesMutex.Lock()
	sci.stubsMutex.Lock()

	defer sci.stubsMutex.Unlock()
	defer sci.servicesMutex.Unlock()
	_, ok := sci.servicesAddress[serviceName]
	if !ok {
		logs.Error("get rpc stub failed: there is no this %s service client", serviceName)
		return nil, errors.New(fmt.Sprintf("get rpc stub failed: there is no this %s service client", serviceName))
	}
	address := sci.servicesAddress[serviceName]
	_, ok = sci.stubs[serviceName]
	if !ok {
		// first use
		stub, err := NewGrpcStub(address)
		if err != nil {
			return nil, err
		}
		sci.stubs[address] = stub
	}
	return sci.stubs[address], nil
}
