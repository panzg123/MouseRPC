package mouserpc

import (
	"fmt"
	"net"
)

func NewClient() *Client {
	return &Client{}
}

type Client struct{}

func (c *Client) Invoke(target string) error {
	addr, err := net.ResolveUDPAddr("udp", target)
	if err != nil {
		fmt.Println("Can't resolve address: ", err)
		return err
	}
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		fmt.Println("Can't dial: ", err)
		return err
	}
	defer conn.Close()
	_, err = conn.Write([]byte("hell world"))
	if err != nil {
		fmt.Println("failed:", err)
		return err
	}
	data := make([]byte, 1024)
	n, err := conn.Read(data)
	if err != nil {
		fmt.Println("failed to read UDP msg because of ", err)
		return err
	}
	fmt.Println("receive data succ ", string(data[0:n]))
	return nil
}
