package server

import (
	"context"
	"fmt"
	"github.com/ZuoFuhong/grpc-cgi-proxy/consts"
	"github.com/ZuoFuhong/grpc-cgi-proxy/pkg/log"
	codec "github.com/ZuoFuhong/grpc-middleware/encoding/json"
	meta "github.com/ZuoFuhong/grpc-middleware/metadata"
	gm "github.com/ZuoFuhong/grpc-middleware/tracing"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"time"
)

type Handler struct {
	serviceName string
	namespace   string
	cmd         string
	timeout     int
}

func NewHandler(serviceName, namespace, cmd string, timeout int) *Handler {
	return &Handler{
		serviceName: serviceName,
		namespace:   namespace,
		cmd:         cmd,
		timeout:     timeout,
	}
}

func (h *Handler) handle(ctx context.Context, req interface{}) (map[string]interface{}, error) {
	// 使用 monica 注册中心
	target := fmt.Sprintf("monica://%s/%s", h.namespace, h.serviceName)
	method := fmt.Sprintf("/%s.%s/%s", h.serviceName, h.serviceName, h.cmd)
	conn, err := newConn(ctx, target)
	if err != nil {
		return nil, err
	}
	// 注入元数据
	traceId := ctx.Value(consts.TraceId).(string)
	ctx = metadata.NewIncomingContext(ctx, metadata.New(map[string]string{meta.TraceId: traceId}))
	// 超时控制
	ctx, cancel := context.WithTimeout(ctx, time.Second*time.Duration(h.timeout))
	defer cancel()
	rsp := make(map[string]interface{})
	if err := conn.Invoke(ctx, method, req, &rsp, grpc.CallContentSubtype(codec.Name)); err != nil {
		log.ErrorContextf(ctx, "grpc.Invoke failed, err: %v", err)
		return nil, err
	}
	return rsp, nil
}

// 后续：复用连接优化性能
func newConn(ctx context.Context, target string) (*grpc.ClientConn, error) {
	conn, err := grpc.DialContext(ctx, target, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithUnaryInterceptor(gm.UnaryClientInterceptor()))
	if err != nil {
		log.ErrorContextf(ctx, "grpc.DialContext failed, err: %v", err)
		return nil, err
	}
	return conn, nil
}
