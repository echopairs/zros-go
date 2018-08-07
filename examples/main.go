package main

import (
    "google.golang.org/grpc"
    "reflect"
    "errors"
)

const (
    address = "localhost:50051"
)

func main() {
    grpc.Dial(address, grpc.WithInsecure())
}

func extractElement(v reflect.Value) (interface{}, error) {
    if v.Kind() != reflect.Ptr {
        return nil, errors.New("invalid input")
    }

    v = v.Elem()
    var elem interface{}
    elem = v.Interface()
    return elem, nil
}