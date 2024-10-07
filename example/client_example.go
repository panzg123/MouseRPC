package main

import (
	"context"
	"fmt"

	helloworld "github.com/panzg123/mouserpc/stub"
)

func main() {
	req := &helloworld.HelloRequest{Name: "hello world!!!"}
	cli := helloworld.NewClientProxy("127.0.0.1:9090")
	rsp, err := cli.SayHi(context.Background(), req)
	if err != nil {
		fmt.Println("client invoke failed, err =  ", err)
		return
	}
	fmt.Printf("sayHi success, req: %+v, rsp: %+v\n", req, rsp)
}
