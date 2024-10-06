package main

import (
	"context"

	"github.com/panzg123/mouserpc"
	pb "github.com/panzg123/mouserpc/stub"
)

// HiService 接口实现
type HiService struct {
}

// SayHi ...
func (h *HiService) SayHi(ctx context.Context, req *pb.HelloRequest) (*pb.HelloReply, error) {
	rsp := &pb.HelloReply{
		Name: req.GetName(),
	}
	return rsp, nil
}

func main() {
	s := mouserpc.NewServer()
	pb.RegisterServer(s, &HiService{})
	_ = s.Serve()
}
