package main

import (
	"fmt"

	"github.com/panzg123/mouserpc"
)

func main() {
	c := mouserpc.NewClient()
	if err := c.Invoke("127.0.0.1:9090"); err != nil {
		fmt.Println("client invoke failed, err =  ", err)
		return
	}
}
