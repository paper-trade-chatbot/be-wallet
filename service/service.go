package service

import (
	"context"
	"fmt"

	"github.com/paper-trade-chatbot/be-wallet/logging"

	memberGrpc "github.com/paper-trade-chatbot/be-proto/member"
	"github.com/paper-trade-chatbot/be-wallet/config"
	"github.com/paper-trade-chatbot/be-wallet/service/member"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

var Impl ServiceImpl
var (
	MemberServiceHost    = config.GetString("MEMBER_GRPC_HOST")
	MemberServerGRpcPort = config.GetString("MEMBER_GRPC_PORT")
	memberServiceConn    *grpc.ClientConn
)

type ServiceImpl struct {
	MemberIntf member.MemberIntf
}

func GrpcDial(addr string) (*grpc.ClientConn, error) {
	return grpc.Dial(addr, grpc.WithInsecure(), grpc.WithDefaultCallOptions(
		grpc.MaxCallRecvMsgSize(20*1024*1024),
		grpc.MaxCallSendMsgSize(20*1024*1024)), grpc.WithUnaryInterceptor(clientInterceptor))
}

func Initialize(ctx context.Context) {

	var err error

	addr := MemberServiceHost + ":" + MemberServerGRpcPort
	fmt.Println("dial to order grpc server...", addr)
	memberServiceConn, err = GrpcDial(addr)
	if err != nil {
		fmt.Println("Can not connect to gRPC server:", err)
	}
	fmt.Println("dial done")
	memberConn := memberGrpc.NewMemberServiceClient(memberServiceConn)
	Impl.MemberIntf = member.New(memberConn)
}

func Finalize(ctx context.Context) {
	memberServiceConn.Close()
}

func clientInterceptor(ctx context.Context, method string, req interface{}, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	requestId, _ := ctx.Value(logging.ContextKeyRequestId).(string)
	account, _ := ctx.Value(logging.ContextKeyAccount).(string)

	ctx = metadata.NewOutgoingContext(ctx, metadata.New(map[string]string{
		logging.ContextKeyRequestId: requestId,
		logging.ContextKeyAccount:   account,
	}))

	err := invoker(ctx, method, req, reply, cc, opts...)
	if err != nil {
		fmt.Println("clientInterceptor err:", err.Error())
	}

	return err
}
