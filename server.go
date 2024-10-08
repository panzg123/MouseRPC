package mouserpc

import (
	"context"
	"fmt"
	"net"
	"os"

	"github.com/panzg123/mouserpc/codec"
	"google.golang.org/protobuf/proto"
)

// DefaultRecvLen 默认收包大小
var DefaultRecvLen = 65536

type RPCHandler func(ctx context.Context, req []byte) (interface{}, error)

// Server represents an RPC Server.
type Server struct {
	handlers map[string]RPCHandler
}

func (s *Server) RegisterHandler(rpcName string, h RPCHandler) {
	if s.handlers == nil {
		s.handlers = make(map[string]RPCHandler)
	}
	s.handlers[rpcName] = h
}

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
		data := make([]byte, DefaultRecvLen)
		size, remoteAddr, err := conn.ReadFromUDP(data)
		if err != nil {
			fmt.Println("failed to read UDP msg because of ", err.Error())
			continue
		}
		// 检查包的完整性，不完整则抛弃
		msg, err := codec.DefaultServerCodec.Decode(data[0:size])
		if err != nil {
			fmt.Printf("ReadMsg failed, err[%v]\n", err)
			continue
		}
		// todo decode msg，传入闭包函数，解码req
		// 处理逻辑
		fmt.Printf("req header[%+v]\n", msg.ReqHeader)
		fmt.Printf("recv size[%d] data.len[%d] remote addr[%v]\n", size, len(data), remoteAddr)
		h, ok := s.handlers[msg.ReqHeader.InterfaceName]
		if !ok {
			// TODO 处理失败，要回包错误
			fmt.Printf("interface name not registed, name = %s\n", msg.ReqHeader.InterfaceName)
			continue
		}
		// todo ctx控制超时等信息，此处reqBody还是[]byte类型
		ctx := context.Background()
		rsp, err := h(ctx, msg.ReqBody)
		// TODO Encode，此处rspBody是proto.Message类型
		msg.RspBody, msg.Err = proto.Marshal(rsp.(proto.Message))
		if msg.Err != nil {
			fmt.Printf("marshal failed, err: %v\n", err)
			continue
		}
		fmt.Printf("handle msg success, rsp: %+v, err: %v\n", rsp, err)
		// 编码回包
		rspData, err := codec.DefaultServerCodec.Encode(msg)
		if err != nil {
			fmt.Printf("encode failed, err: %v\n", err)
			continue
		}
		// 回包
		size, err = conn.WriteToUDP(rspData, remoteAddr)
		if err != nil {
			fmt.Println("write response failed, err = ", err)
			return err
		}
		fmt.Printf("write response success, rspData.size: %d, write.size: %d\n", len(rspData), size)
	}
}
