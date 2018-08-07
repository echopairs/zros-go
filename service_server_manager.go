package zros_go

import (
	"sync"
	"context"
	"net"

	"github.com/astaxie/beego/logs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb "zros-go/zros_rpc"
)

var  _ ServiceServerManager = (*GrpcServerImpl)(nil)

type ServiceServerManager interface {
	Start() (string, error)
	Stop()
	RegisterServer(server *ServiceServer) (error)
	UnregisterServer(serviceName string) (error)
}

type GrpcServerImpl struct {
	ServiceAddress string
	mutex 		  sync.RWMutex
	servers 	  map[string]*ServiceServer
	lis  		  net.Listener
}

func NewGrpcServerImpl() (*GrpcServerImpl){
	return &GrpcServerImpl{
		servers: make(map[string]*ServiceServer),
	}
}

func (gsi *GrpcServerImpl) Start() (string, error) {
	lis, err := net.Listen("tcp", "127.0.0.1:")
	if err != nil {
		logs.Error("GrpcServerImpl Start Listen error %v", err)
		return "", err
	}
	gsi.ServiceAddress = lis.Addr().String()
	gsi.lis = lis
	go gsi.serve()
	return gsi.ServiceAddress, nil
}

func (gsi *GrpcServerImpl) serve() {
	s := grpc.NewServer()
	pb.RegisterServiceRPCServer(s, gsi)

	// Register reflection service on gRPC server.
	reflection.Register(s)
	if err := s.Serve(gsi.lis); err != nil {
		logs.Error("failed to serve: %v", err)
	}
}


func (gsi *GrpcServerImpl) Stop() {

}

func (gsi *GrpcServerImpl) RegisterServer(server *ServiceServer) (error) {
	// 1. register to master first
	gsd := GetGlobalServiceDiscovery()
	err := gsd.AddServiceServer(server)
	if err != nil {
		logs.Error("AddServiceServer to master failed")
		return err
	}
	// 2. register to memory
	gsi.mutex.Lock()
	defer gsi.mutex.Unlock()
	serviceName := server.GetServiceName()
	_, ok := gsi.servers[serviceName]
	if ok {
		logs.Warn("%s already register", serviceName)
		return nil
	}
	gsi.servers[serviceName] = server
	logs.Info("RegisterServer %s ok", serviceName)
	return nil
}

func (gsi *GrpcServerImpl) UnregisterServer(serviceName string) (error) {
	gsi.mutex.Lock()
	defer gsi.mutex.Unlock()
	_, ok := gsi.servers[serviceName]
	if ok {
		delete(gsi.servers, serviceName)
	}
	return nil
}

func (gsi *GrpcServerImpl) InvokeService(ctx context.Context, request *pb.ServiceRequest) (*pb.ServiceResponse, error) {
	serviceName := request.GetServiceName()
	gsi.mutex.Lock()
	defer gsi.mutex.Unlock()
	server, ok := gsi.servers[serviceName]
	if !ok {
		res := &pb.ServiceResponse{}
		status := &pb.Status{}
		status.Code = pb.Status_FAILED_PRECONDITION
		status.Details = "There is no server for this service " + serviceName
		res.Status = status
		return res, nil
	}
	return server.Invoke(request)
}