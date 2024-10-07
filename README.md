# MouseRPC

## 方案设计

基于protobuf3，定义简单的req rsp

### RPC协议格式

frame header + rpc header + req  或者 frame header + rpc header + rsp

frame固定6字节： 2 byte magic number + 2 byte rpc header len + 2 byte total header

magic number: 0x1024

rpc header见proto/rpc.proto

### RPC流图

client --> filter --> encode --> transport --> decode --> filter --> server

## 开发计划

### step1 只支持简单的单发单收模式，无service&interface桩代码
- [ ] server
- [ ] client

### step2 支持pb定义桩代码，生成client server桩代码
- [ ] codec
- [ ] stub
- [ ] filter

### step3 服务注册与发现

### step4 负载均衡

### step5 支持http协议