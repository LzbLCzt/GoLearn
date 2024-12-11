package gRPCStream

import (
	"GoLearn/web/chapter5/proto"
	"context"
	"fmt"
	"io"
)

type HelloServiceImpl struct{}

func (p *HelloServiceImpl) Hello(
	ctx context.Context, args *proto.String,
) (*proto.String, error) {
	reply := &proto.String{Value: "hello:" + args.GetValue()}
	return reply, nil
}

//todo 实现流服务
/*
todo 服务端在循环中接收客户端发来的数据，如果遇到io.EOF表示客户端流被关闭，
todo 如果函数退出表示服务端流关闭。生成返回的数据通过流发送给客户端，双向流数据的发送和接收都是完全独立的行为
 */
func (p *HelloServiceImpl) Channel(stream proto.HelloService_ChannelServer) error {
	for {
		args, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
		msg := fmt.Sprintf("received you name: %v" + args.GetValue())
		fmt.Println(msg)
		reply := &proto.String{Value: msg}

		err = stream.Send(reply)
		if err != nil {
			return err
		}
	}
}