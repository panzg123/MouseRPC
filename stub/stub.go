package helloworld

import (
	"context"

	"github.com/panzg123/mouserpc"
	"google.golang.org/protobuf/proto"
)

// Service 定义一个服务的对外接口
// todo 工具生成桩代码
type Service interface {
	// SayHi ...
	SayHi(ctx context.Context, req *HelloRequest) (*HelloReply, error)
}

// Service_SayHi_Handler handler桩代码生成
func Service_SayHi_Handler(svr interface{}, ctx context.Context, reqBody []byte) (interface{}, error) {
	req := &HelloRequest{}
	// todo 传入匿名函数 + codec
	if err := proto.Unmarshal(reqBody, req); err != nil {
		return nil, err
	}
	return svr.(Service).SayHi(ctx, req)
}

// RegisterServer 注册service
// todo 拆分到service
func RegisterServer(s *mouserpc.Server, imp Service) {
	h := func(ctx context.Context, req []byte) (interface{}, error) {
		return Service_SayHi_Handler(imp, ctx, req)
	}
	s.RegisterHandler("SayHi", h)
}
