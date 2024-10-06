package mouserpc

import (
	"encoding/binary"
	"fmt"
	"log"
	"math/rand"
	"net"

	"github.com/panzg123/mouserpc/rpcproto"
	"google.golang.org/protobuf/proto"
)

func NewClient() *Client {
	return &Client{}
}

type Client struct{}

// Invoke ...
// todo rpcName 放在context.msg中
func (c *Client) Invoke(target string, rpcName string, req interface{}, rsp interface{}) error {
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
	// encode
	reqData, err := encode(rpcName, req)
	if err != nil {
		log.Printf("client encode failed, err: %v\n", err)
		return err
	}
	_, err = conn.Write(reqData)
	if err != nil {
		fmt.Println("failed:", err)
		return err
	}
	data := make([]byte, 1024)
	readSize, err := conn.Read(data)
	if err != nil {
		fmt.Println("failed to read UDP msg because of ", err)
		return err
	}
	// todo decode
	fmt.Println("receive data success, size=", string(data[0:readSize]))
	if err = decode(data[:readSize], rsp); err != nil {
		fmt.Printf("decode failed, recvSize: %d, err: %v", readSize, err)
		return err
	}
	fmt.Printf("decode success, rsp: %v\n", rsp)
	return nil
}

func encode(rpcName string, req interface{}) ([]byte, error) {
	header := &rpcproto.RequestHeader{
		AppName:       "",
		ServiceName:   "",
		InterfaceName: rpcName,
		RequestId:     uint32(rand.Int31()),
	}
	headerBuf, err := proto.Marshal(header)
	if err != nil {
		return nil, err
	}
	bodyBuf, err := proto.Marshal(req.(proto.Message))
	if err != nil {
		return nil, err
	}
	totalLen := 6 + len(headerBuf) + len(bodyBuf)
	reqData := make([]byte, totalLen)
	binary.BigEndian.PutUint16(reqData[:2], FrameHeaderMagic)
	binary.BigEndian.PutUint16(reqData[2:4], uint16(len(headerBuf)))
	binary.BigEndian.PutUint16(reqData[4:6], uint16(totalLen))
	copy(reqData[6:6+len(headerBuf)], headerBuf)
	copy(reqData[6+len(headerBuf):], bodyBuf)
	return reqData, nil
}

func decode(rspData []byte, rsp interface{}) error {
	// todo 解析
	return nil
}
