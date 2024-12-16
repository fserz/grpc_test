package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	pb "grpc_test/pb"
	"io"
	"log"
	"net"
	"strings"
	"time"
)

type server struct {
	pb.UnimplementedGreeterServer
}

func (s *server) UnarySayHelloSay(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
	// 加入metadata
	// 通过defer中设置trailer.
	defer func() {
		trailer := metadata.Pairs("timestamp", time.Now().String())
		grpc.SetTrailer(ctx, trailer)
	}()

	// 从客户端请求上下文中读取metadata.
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.DataLoss, "UnarySayHello: failed to get metadata")
	}
	if t, ok := md["token"]; ok {
		fmt.Printf("token from metadata: %+v\n", t)
		if len(t) < 1 || t[0] != "app-test-q1mi" {
			return nil, status.Error(codes.Unauthenticated, "认证失败")
		}
	}

	// 创建和发送header.
	header := metadata.New(map[string]string{"location": "BeiJing"})
	grpc.SendHeader(ctx, header)

	fmt.Printf("request received: %v, say hello...\n", req)

	return &pb.HelloResponse{Reply: req.Name}, nil
}

func (s *server) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
	// 普通RPC
	return &pb.HelloResponse{Reply: "Hello ball ball~"}, nil
}

func (s *server) BidiHello(stream pb.Greeter_BidiHelloServer) error {
	for {
		// 接受流式响应
		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		reply := magic(in.GetName())

		// 发送流式请求
		if err = stream.Send(&pb.HelloResponse{Reply: reply}); err != nil {
			log.Fatalf("send stream reply err: %+v", err)
			return err
		}

	}
}

func magic(s string) string {
	s = strings.ReplaceAll(s, "吗", "")
	s = strings.ReplaceAll(s, "吧", "")
	s = strings.ReplaceAll(s, "你", "我")
	s = strings.ReplaceAll(s, "？", "!")
	s = strings.ReplaceAll(s, "?", "!")
	return s
}

func main() {
	lis, err := net.Listen("tcp", ":8972")
	if err != nil {
		fmt.Printf("failed to listen: %+v", lis)
		return
	}
	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &server{})
	err = s.Serve(lis)
	if err != nil {
		fmt.Printf("failed to server %+v", err)
		return
	}

}
