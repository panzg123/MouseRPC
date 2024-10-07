package mouserpc

import (
	"fmt"
	"log"
	"net"

	"github.com/panzg123/mouserpc/codec"
)

func NewClient() *Client {
	return &Client{}
}

type Client struct{}

// Invoke ...
// todo rpcName 放在context.msg中
func (c *Client) Invoke(target string, msg *codec.Msg) (*codec.Msg, error) {
	addr, err := net.ResolveUDPAddr("udp", target)
	if err != nil {
		fmt.Println("Can't resolve address: ", err)
		return nil, err
	}
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		fmt.Println("Can't dial: ", err)
		return nil, err
	}
	defer conn.Close()
	// encode
	reqData, err := codec.DefaultClientCodec.Encode(msg)
	if err != nil {
		log.Printf("client encode failed, err: %v\n", err)
		return nil, err
	}
	_, err = conn.Write(reqData)
	if err != nil {
		fmt.Println("failed:", err)
		return nil, err
	}
	data := make([]byte, 1024)
	readSize, err := conn.Read(data)
	if err != nil {
		fmt.Println("failed to read UDP msg because of ", err)
		return nil, err
	}
	// decode
	fmt.Println("receive data success, size=", readSize)
	replyMsg, err := codec.DefaultClientCodec.Decode(data[:readSize])
	if err != nil {
		fmt.Printf("decode failed, recvSize: %d, err: %v", readSize, err)
		return nil, err
	}
	fmt.Printf("decode success,rsp.header: %+v, rsp.body.size: %d\n", replyMsg.RspHeader, len(replyMsg.RspBody))
	// 匹配请求id
	if msg.ReqHeader.RequestId != replyMsg.RspHeader.RequestId {
		fmt.Printf("requestId not match, reqId: %d, rspId: %d\n",
			msg.ReqHeader.RequestId, replyMsg.RspHeader.RequestId)
		return replyMsg, fmt.Errorf("requestId not match, reqId: %d, rspId: %d",
			msg.ReqHeader.RequestId, replyMsg.RspHeader.RequestId)
	}
	return replyMsg, nil
}
