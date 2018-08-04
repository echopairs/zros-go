package zros_go

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


type Message interface {

}

type Publisher interface {
	Publisher(msg Message)
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
	AdvertiseService(service string, reqType Message, resType Message, callback interface{}) (ServiceServer, error)
	ServiceClient(service string, reqType Message, resType Message) (ServiceClient, error)
	Advertise(topic string, msgType Message) (Publisher, error)
	Subscriber(topic string, msgType Message, callback interface{}) (Subscriber, error)
}

