package mouserpc

import (
	"fmt"
	"testing"

	pb "github.com/panzg123/mouserpc/rpcproto"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"
)

func TestMarshal(t *testing.T) {
	req := &pb.RequestHeader{
		RequestId: 1,
	}
	buf, err := proto.Marshal(req)
	if err != nil {
		fmt.Printf("marshal failed\n")
		return
	}
	fmt.Printf("marshal success, len[%d]\n", len(buf))
	assert.NotZero(t, len(buf))
}
