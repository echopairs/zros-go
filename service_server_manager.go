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



type ServiceServerManager interface {
	Start() error
	Stop()
	RegisterServer(server *IServiceServer) (error)
	UnregisterServer(serviceName string) (error)
}

type GrpcServerImpl struct {
	ServiceAddress string
	mutex 		  sync.RWMutex
	servers 	  map[string]*ServiceServer
	lis  		  net.Listener
}

func (gsi *GrpcServerImpl) Start() error {
	lis, err := net.Listen("tcp", "127.0.0.1:")
	if err != nil {
		logs.Error("GrpcServerImpl Start Listen error %v", err)
		return err
	}
	gsi.ServiceAddress = lis.Addr().String()
	gsd.lis = lis
	go gsi.serve()
	return nil
}

func (gsi *GrpcServerImpl) serve() {
	s := grpc.NewServer()
	pb.RegisterServiceRPCServer(s, &GrpcServerImpl{})

	// Register reflection service on gRPC server.
	reflection.Register(s)
	if err := s.Serve(gsi.lis); err != nil {
		logs.Error("failed to serve: %v", err)
	}
}


func (gsi *GrpcServerImpl) Stop() {

}

func (gsi *GrpcServerImpl) RegisterServer(server *ServiceServer) (error) {
	gsi.mutex.Lock()
	defer gsi.mutex.Unlock()
	serviceName := server.GetServiceName()
	_, ok := gsi.servers[serviceName]
	if ok {
		logs.Warn("%s already register", serviceName)
	}
	gsi.servers[serviceName] = server
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
	logs.Info("receiver rpc call %s ", serviceName)
	gsi.mutex.Lock()
	defer gsi.mutex.Unlock()
	server, ok := gsi.servers[serviceName]
	if !ok {
		res := &pb.ServiceResponse{}
		res.Status.Code = pb.Status_FAILED_PRECONDITION
		res.Status.Details = "There is no server for this service " + serviceName
		return res, nil
	}
	return server.Invoke(request)
}