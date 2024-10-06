package mouserpc

import (
	"context"
	"encoding/binary"
	"fmt"
	"net"
	"os"

	"github.com/panzg123/mouserpc/rpcproto"
	"google.golang.org/protobuf/proto"
)

// DefaultRecvLen 默认收包大小
var DefaultRecvLen = 65536
var FrameHeaderMagic uint16 = 0x1024

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
		msg, err := ReadMsg(data[0:size])
		if err != nil {
			fmt.Printf("ReadMsg failed, err[%v]\n", err)
			continue
		}
		// todo decode msg，传入闭包函数，解码req
		// 处理逻辑
		fmt.Printf("req header[%+v]\n", msg.reqHeader)
		fmt.Printf("recv size[%d] data.len[%d] remote addr[%v]\n", size, len(data), remoteAddr)
		h, ok := s.handlers[msg.reqHeader.InterfaceName]
		if !ok {
			// TODO 处理失败，要回包错误
			fmt.Printf("interface name not registed, name = %s\n", msg.reqHeader.InterfaceName)
			continue
		}
		// todo ctx控制超时等信息，此处reqBody还是[]byte类型
		ctx := context.Background()
		rsp, err := h(ctx, msg.reqBody)
		// TODO Encode，此处rspBody是proto.Message类型
		rspBody, err := proto.Marshal(rsp.(proto.Message))
		if err != nil {
			fmt.Printf("marshal failed, err: %v\n", err)
			continue
		}
		fmt.Printf("handle msg success, rsp: %+v, err: %v\n", rsp, err)
		// 编码回包
		rspData := WriteMsg(msg, err, rspBody)
		// 回包
		size, err = conn.WriteToUDP(rspData, remoteAddr)
		if err != nil {
			fmt.Println("write response failed, err = ", err)
			return err
		}
		fmt.Printf("write response success, rspData.size: %d, write.size: %d\n", len(rspData), size)
	}
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
	headLen := binary.BigEndian.Uint16(buf[2:4])
	totalLen := binary.BigEndian.Uint16(buf[4:6])
	if totalLen != uint16(len(buf)) {
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

// WriteMsg 写入回包信息
func WriteMsg(msg *Msg, err error, rspBody []byte) []byte {
	header := &rpcproto.ResponseHeader{
		AppName:       msg.reqHeader.AppName,
		ServiceName:   msg.reqHeader.ServiceName,
		InterfaceName: msg.reqHeader.InterfaceName,
		RequestId:     msg.reqHeader.RequestId,
		Ret:           0, // todo 错误码赋值
		Msg:           "",
	}
	headerBuf, err := proto.Marshal(header)
	if err != nil {
		fmt.Printf("writeMsg marshal header failed, header: %+v. err: %v", header, err)
		return nil
	}
	headLen := len(headerBuf)
	rspBodyLen := len(rspBody)
	totalLen := 6 + headLen + rspBodyLen
	buf := make([]byte, totalLen)
	binary.BigEndian.PutUint16(buf[:2], uint16(FrameHeaderMagic))
	binary.BigEndian.PutUint16(buf[2:4], uint16(headLen))
	binary.BigEndian.PutUint16(buf[4:6], uint16(totalLen))
	copy(buf[6:6+headLen], headerBuf)
	copy(buf[6+headLen:], rspBody)
	return buf
}
