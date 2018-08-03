package main

import (
    "google.golang.org/grpc"
)

const (
    address = "localhost:50051"
)

func main() {
    grpc.Dial(address, grpc.WithInsecure())
}
