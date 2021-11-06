package mouserpc

import (
	"encoding/binary"
	"fmt"
	"net"
	"os"

	"github.com/panzg123/mouserpc/rpcproto"
	"google.golang.org/protobuf/proto"
)

// DefaultRecvLen 默认收包大小
var DefaultRecvLen = 65536
var FrameHeaderMagic = 0x1024

type RPCHandler func(req proto.Message, rsp proto.Message) error

// Server represents an RPC Server.
type Server struct {
	handlers map[string]RPCHandler
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
		msg, err := ReadMsg(data[0:size])
		if err != nil {
			fmt.Printf("ReadMsg failed, err[%v]\n", err)
			continue
		}
		// 处理逻辑
		fmt.Printf("req header[%+v]\n", msg.reqHeader)
		h, ok := s.handlers[msg.reqHeader.InterfaceName]
		if !ok {
			// TODO 处理失败，要回包错误
			fmt.Printf("interface name not registed, name = %s", msg.reqHeader.InterfaceName)
			continue
		}
		// TODO 这里如何确认是哪个类型呢

		// TODO Encode

		// 写回包
		fmt.Printf("recv size[%d] data.len[%d] remote addr[%v]\n", size, len(data), remoteAddr)
		// 回包
		size, err = conn.WriteToUDP(data[0:size], remoteAddr)
		if err != nil {
			fmt.Println("write response failed, err = ", err)
			return err
		}
		fmt.Println("write response success, n = ", size)
	}
	return nil
}

// Msg 一条完整的rpc消息
type Msg struct {
	reqHeader *rpcproto.RequestHeader
	rspHeader *rpcproto.ResponseHeader
	reqBody   []byte
	rspBody   []byte
}

func ReadMsg(buf []byte) (*Msg, error) {
	m := &Msg{}
	// 判断包的完整性
	if len(buf) < 6 {
		return nil, fmt.Errorf("buf len invalid, len = %d", len(buf))
	}
	magic := binary.BigEndian.Uint16(buf[:2])
	if magic != uint16(FrameHeaderMagic) {
		return nil, fmt.Errorf("frame invalid, magic is %d", magic)
	}
	// 解析出header
	headLen := binary.BigEndian.Uint32(buf[2:4])
	totalLen := binary.BigEndian.Uint32(buf[4:6])
	if totalLen != uint32(len(buf)) {
		return nil, fmt.Errorf("total len invalid")
	}
	h := &rpcproto.RequestHeader{}
	if err := proto.Unmarshal(buf[6:6+headLen], h); err != nil {
		return nil, err
	}
	m.reqHeader = h
	// 解析出body
	m.reqBody = buf[6+headLen : totalLen]
	return m, nil
}
