package zros_go

import "testing"

var masterAddress = "localhost:23333"

func TestNewGrpcServiceDiscovery(t *testing.T) {
	_, err := NewGrpcServiceDiscovery(masterAddress)
	if err != nil {
		t.Error(err)
	}
}

func TestGrpcServiceDiscovery_IsConnectedToMaster(t *testing.T) {
	gsd, _ := NewGrpcServiceDiscovery(masterAddress)
	err := gsd.IsConnectedToMaster()
	if err != nil {
		t.Error(err)
	}
}

func TestGrpcServiceDiscovery_Spin(t *testing.T) {
	gsd, _:= NewGrpcServiceDiscovery(masterAddress)
	gsd.Spin()
}

//func TestGrpcServiceDiscovery_AddServiceServer(t *testing.T) {
//	gsd, _ := NewGrpcServiceDiscovery(masterAddress)
//	server := Ne
//}