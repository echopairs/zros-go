package zros_go


import (
	pb "zros-go/zros_rpc"
	"google.golang.org/grpc"
	"context"
	"time"
	"net"
	"google.golang.org/grpc/reflection"
	"github.com/astaxie/beego/logs"
)

//type ServiceDiscovery interface {
//	AddServiceServer() error
//	AddServiceClient() error
//	AddPublisher()	error
//	AddSubscriber() error
//}


var stubCallShortTimeOut = 5*1000*1000
var stubCallLongTimeOut = 3000*1000*1000


type GrpcServiceDiscovery struct {
	masterRpcStub 	pb.MasterRPCClient
	conn 			*grpc.ClientConn
	agentAddress	string
	lis 			net.Listener
}

func NewGrpcServiceDiscovery(address string) (*GrpcServiceDiscovery, error) {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	c := pb.NewMasterRPCClient(conn)
	gsd := &GrpcServiceDiscovery{
		masterRpcStub:c,
		conn:conn,
	}
	return gsd, nil
}

func (gsd *GrpcServiceDiscovery) Spin() error {
	lis, err := net.Listen("tcp", "127.0.0.1:")
	if err != nil {
		logs.Error("Grpc Service Discovery Spin Listen error %v", err)
		return err
	}
	logs.Info("the address is %s", lis.Addr().String())
	gsd.agentAddress = lis.Addr().String()
	gsd.lis = lis
	go gsd.serve()
	return nil
}

func (gsd *GrpcServiceDiscovery) InvokeService(context.Context, *pb.ServiceRequest) (*pb.ServiceResponse, error) {
	return nil, nil
}

func (gsd *GrpcServiceDiscovery) serve() {
	s := grpc.NewServer()
	pb.RegisterServiceRPCServer(s, &GrpcServiceDiscovery{})

	// Register reflection service on gRPC server.
	reflection.Register(s)
	if err := s.Serve(gsd.lis); err != nil {
		logs.Error("failed to serve: %v", err)
	}
}

func (gsd *GrpcServiceDiscovery) IsConnectedToMaster() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(stubCallShortTimeOut))
	defer cancel()
	in := &pb.PingRequest{}
	status, err := gsd.masterRpcStub.Ping(ctx, in)
	if err != nil {
		return err
	}
	if status.Code != pb.Status_OK {
		return err
	}
	return nil
}


