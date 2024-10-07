package codec

import (
	"encoding/binary"
	"fmt"

	"github.com/panzg123/mouserpc/rpcproto"
	"google.golang.org/protobuf/proto"
)

// StandardServerCodec 实现一个标准的解码器，见ReadMe协议规范，服务端用
type StandardServerCodec struct {
}

// StandardClientCodec 实现一个标准的编码器，见ReadMe协议规范，客户端用
type StandardClientCodec struct {
}

func (s *StandardClientCodec) Encode(msg *Msg) ([]byte, error) {
	headerBuf, err := proto.Marshal(msg.ReqHeader)
	if err != nil {
		return nil, err
	}
	totalLen := 6 + len(headerBuf) + len(msg.ReqBody)
	reqData := make([]byte, totalLen)
	binary.BigEndian.PutUint16(reqData[:2], FrameHeaderMagic)
	binary.BigEndian.PutUint16(reqData[2:4], uint16(len(headerBuf)))
	binary.BigEndian.PutUint16(reqData[4:6], uint16(totalLen))
	copy(reqData[6:6+len(headerBuf)], headerBuf)
	copy(reqData[6+len(headerBuf):], msg.ReqBody)
	return reqData, nil
}

func (s *StandardClientCodec) Decode(buf []byte) (*Msg, error) {
	m := &Msg{}
	// 判断包的完整性
	if len(buf) < 6 {
		return nil, fmt.Errorf("buf len invalid, len = %d", len(buf))
	}
	magic := binary.BigEndian.Uint16(buf[:2])
	if magic != FrameHeaderMagic {
		return nil, fmt.Errorf("frame invalid, magic is %d", magic)
	}
	// 解析出header
	headLen := binary.BigEndian.Uint16(buf[2:4])
	totalLen := binary.BigEndian.Uint16(buf[4:6])
	if totalLen != uint16(len(buf)) {
		return nil, fmt.Errorf("total len invalid")
	}
	h := &rpcproto.ResponseHeader{}
	if err := proto.Unmarshal(buf[6:6+headLen], h); err != nil {
		return nil, err
	}
	m.RspHeader = h
	// 解析出body
	m.RspBody = buf[6+headLen : totalLen]
	return m, nil
}

func (s *StandardServerCodec) Encode(msg *Msg) ([]byte, error) {
	header := &rpcproto.ResponseHeader{
		AppName:       msg.ReqHeader.AppName,
		ServiceName:   msg.ReqHeader.ServiceName,
		InterfaceName: msg.ReqHeader.InterfaceName,
		RequestId:     msg.ReqHeader.RequestId,
		Ret:           0, // todo 错误码赋值
		Msg:           "",
	}
	headerBuf, err := proto.Marshal(header)
	if err != nil {
		fmt.Printf("writeMsg marshal header failed, header: %+v. err: %v", header, err)
		return nil, err
	}
	headLen := len(headerBuf)
	rspBodyLen := len(msg.RspBody)
	totalLen := 6 + headLen + rspBodyLen
	buf := make([]byte, totalLen)
	binary.BigEndian.PutUint16(buf[:2], FrameHeaderMagic)
	binary.BigEndian.PutUint16(buf[2:4], uint16(headLen))
	binary.BigEndian.PutUint16(buf[4:6], uint16(totalLen))
	copy(buf[6:6+headLen], headerBuf)
	copy(buf[6+headLen:], msg.RspBody)
	return buf, nil
}

func (s *StandardServerCodec) Decode(buf []byte) (*Msg, error) {
	m := &Msg{}
	// 判断包的完整性
	if len(buf) < 6 {
		return nil, fmt.Errorf("buf len invalid, len = %d", len(buf))
	}
	magic := binary.BigEndian.Uint16(buf[:2])
	if magic != FrameHeaderMagic {
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
	m.ReqHeader = h
	// 解析出body
	m.ReqBody = buf[6+headLen : totalLen]
	return m, nil
}
