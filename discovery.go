package zros_go


import (
	"context"
	"time"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"github.com/astaxie/beego/logs"

	pb "zros-go/zros_rpc"
	"fmt"
)


var stubCallShortTimeOut = 5*1000*1000
var stubCallLongTimeOut = 3000*1000*1000



type DiscoveryError struct {
	errorCode int
	detail string
}

func (de *DiscoveryError) Error() string {
	return de.detail
}

type GrpcServiceDiscovery struct {
	masterRpcStub 	pb.MasterRPCClient
	conn 			*grpc.ClientConn
	agentAddress	string
	lis 			net.Listener
}

func NewGrpcServiceDiscovery(masterAddress string) (*GrpcServiceDiscovery, error) {
	conn, err := grpc.Dial(masterAddress, grpc.WithInsecure())
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

func (gsd *GrpcServiceDiscovery) RegisterPublisher(context.Context, *pb.PublisherInfo) (*pb.Status, error) {
	return nil, nil
}

func (gsd *GrpcServiceDiscovery) UnregisterPublisher(context.Context, *pb.PublisherInfo) (*pb.Status, error) {
	return nil, nil
}

func (gsd *GrpcServiceDiscovery) RegisterServiceServer(ctx context.Context, info *pb.ServiceServerInfo) (*pb.Status, error) {
	return nil, nil
}

func (gsd *GrpcServiceDiscovery) UnregisterServiceServer(ctx context.Context, info *pb.ServiceServerInfo) (*pb.Status, error) {
	return nil, nil
}

func (gsd *GrpcServiceDiscovery) Ping(ctx context.Context, request *pb.PingRequest) (*pb.Status, error) {
	return nil, nil
}

func (gsd *GrpcServiceDiscovery) serve() {
	s := grpc.NewServer()
	pb.RegisterServiceDiscoveryRPCServer(s, gsd)

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


// register
func (gsd *GrpcServiceDiscovery) AddServiceServer(server *ServiceServer) (error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(stubCallShortTimeOut))
	logs.Info("AddServiceServer to master")
	defer cancel()
	request := &pb.ServiceServerInfo{}
	physicalNodeInfo := &pb.PhysicalNodeInfo{}

	request.ServiceName = server.GetServiceName()
	physicalNodeInfo.AgentAddress = gsd.agentAddress
	physicalNodeInfo.RealAddress = server.node.NodeAddress
	logs.Info("agentAddress is %s nodeAddress is %s", gsd.agentAddress, server.node.NodeAddress)
	physicalNodeInfo.Name = server.node.NodeName
	request.PhysicalNodeInfo = physicalNodeInfo
	status, err := gsd.masterRpcStub.RegisterServiceServer(ctx, request)

	if err != nil {
		logs.Error(fmt.Sprintf("AddServiceServer %s to master failed for %s", server.GetServiceName(), err.Error()))
		return err
	}

	if status.Code != pb.Status_OK {
		logs.Error(fmt.Sprintf("AddServiceServer %s to master failed for %s", server.GetServiceName(), status.Details))
		return &DiscoveryError{detail:status.Details}

	}
	return nil
}