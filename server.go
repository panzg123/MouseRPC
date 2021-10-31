package mouserpc

import (
	"fmt"
	"net"
	"os"
)

// Server represents an RPC Server.
type Server struct{}

// NewServer returns a new Server.
func NewServer() *Server {
	return &Server{}
}

// Serve 启动所有服务
func (s *Server) Serve() error {
	// 启动一个udp监听服务
	s.ListenAndServer()
	return nil
}

func (s *Server) ListenAndServer() {
	addr, err := net.ResolveUDPAddr("udp", "127.0.0.1:9090")
	if err != nil {
		fmt.Println("Can't resolve address: ", err)
		os.Exit(1)
	}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Println("Error listening:", err)
		os.Exit(1)
	}
	defer conn.Close()
	// 循环处理包
	s.serverPacket(conn)
	// TODO 监听退出信号
}

func (s *Server) serverPacket(conn *net.UDPConn) error {
	for {
		data := make([]byte, 1024)
		n, remoteAddr, err := conn.ReadFromUDP(data)
		if err != nil {
			fmt.Println("failed to read UDP msg because of ", err.Error())
			return err
		}
		fmt.Printf("n[%d] data.len[%d] remote addr[%v]", n, len(data), remoteAddr)
		// 回包
		n, err = conn.WriteToUDP(data[0:n], remoteAddr)
		if err != nil {
			fmt.Println("write response failed, err = ", err)
			return err
		}
		fmt.Println("write response success, n = ", n)
	}
	return nil
}
