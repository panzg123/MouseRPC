// Package codec 编解码器
package codec

import (
	"github.com/panzg123/mouserpc/rpcproto"
)

var FrameHeaderMagic uint16 = 0x1024

// DefaultServerCodec 标准编码，默认编码，后期支持可自定义
var DefaultServerCodec Codec = &StandardServerCodec{}
var DefaultClientCodec Codec = &StandardClientCodec{}

// Codec 编解码
type Codec interface {
	// Encode 编码
	Encode(msg *Msg) (buffer []byte, err error)
	// Decode 解码
	Decode(buffer []byte) (msg *Msg, err error)
}

// Msg 一条标准的rpc消息
type Msg struct {
	ReqHeader *rpcproto.RequestHeader  // 请求头
	RspHeader *rpcproto.ResponseHeader // 响应头
	ReqBody   []byte                   // 请求包体
	RspBody   []byte                   // 响应包体
	Err       error                    // 请求错误信息
}
