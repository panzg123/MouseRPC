// Package helloworld 定义了桩代码，应该由工具生成
package helloworld

import (
	"context"

	"github.com/panzg123/mouserpc"
	"github.com/panzg123/mouserpc/codec"
	"github.com/panzg123/mouserpc/rpcproto"
	"google.golang.org/protobuf/proto"
)

// Service 定义一个服务的对外接口
// todo 工具生成桩代码
type Service interface {
	// SayHi ...
	SayHi(ctx context.Context, req *HelloRequest) (*HelloReply, error)
}

// ServiceSayHiHandler handler桩代码生成
func ServiceSayHiHandler(svr interface{}, ctx context.Context, reqBody []byte) (interface{}, error) {
	req := &HelloRequest{}
	// todo 传入匿名函数 + codec
	if err := proto.Unmarshal(reqBody, req); err != nil {
		return nil, err
	}
	return svr.(Service).SayHi(ctx, req)
}

// RegisterServer 注册service
// todo 桩代码生成，拆分到service中
func RegisterServer(s *mouserpc.Server, imp Service) {
	h := func(ctx context.Context, req []byte) (interface{}, error) {
		return ServiceSayHiHandler(imp, ctx, req)
	}
	s.RegisterHandler("SayHi", h)
}

// NewClientProxy 客户端调用的一个示例
// todo 桩代码生成
var NewClientProxy = func(target string) Service {
	return &ClientImp{
		cli:    mouserpc.NewClient(),
		target: target,
	}
}

type ClientImp struct {
	cli    *mouserpc.Client
	target string
}

func (c *ClientImp) SayHi(ctx context.Context, req *HelloRequest) (*HelloReply, error) {
	// 请求头
	header := &rpcproto.RequestHeader{
		AppName:       "app",
		ServiceName:   "service",
		InterfaceName: "SayHi",
		RequestId:     0,
	}
	// 序列化
	reqBody, err := proto.Marshal(req)
	if err != nil {
		return nil, err
	}
	rspMsg, err := c.cli.Invoke(c.target, &codec.Msg{ReqHeader: header, ReqBody: reqBody})
	if err != nil {
		return nil, err
	}
	// 反序列化
	rsp := &HelloReply{}
	if err = proto.Unmarshal(rspMsg.RspBody, rsp); err != nil {
		return nil, err
	}
	return rsp, nil
}
