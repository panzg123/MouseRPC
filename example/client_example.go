package main

import (
	"fmt"

	"github.com/panzg123/mouserpc"
	helloworld "github.com/panzg123/mouserpc/stub"
)

func main() {
	c := mouserpc.NewClient()
	req := &helloworld.HelloRequest{Name: "hello world!!!"}
	rsp := &helloworld.HelloReply{}
	if err := c.Invoke("127.0.0.1:9090", "SayHi", req, rsp); err != nil {
		fmt.Println("client invoke failed, err =  ", err)
		return
	}
}
