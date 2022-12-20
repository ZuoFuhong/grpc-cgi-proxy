package log

import (
	"context"
	"fmt"
	"github.com/ZuoFuhong/grpc-cgi-proxy/consts"
	"log"
	"os"
)

func Debugf(format string, v ...interface{}) {
	_ = log.Output(2, fmt.Sprintf("DEBUG "+format+"\n", v...))
}

func Infof(format string, v ...interface{}) {
	_ = log.Output(2, fmt.Sprintf("INFO  "+format+"\n", v...))
}

func Warnf(format string, v ...interface{}) {
	_ = log.Output(2, fmt.Sprintf("WARN  "+format+"\n", v...))
}

func Errorf(format string, v ...interface{}) {
	_ = log.Output(2, fmt.Sprintf("ERROR "+format+"\n", v...))
}

func Fatal(v ...interface{}) {
	_ = log.Output(2, fmt.Sprint(v...))
	os.Exit(1)
}

func DebugContextf(ctx context.Context, format string, args ...interface{}) {
	traceId := ctx.Value(consts.TraceId).(string)
	prefix := fmt.Sprintf("DEBUG %s ", traceId)
	_ = log.Output(2, fmt.Sprintf(prefix+format+"\n", args...))
}

func InfoContextf(ctx context.Context, format string, args ...interface{}) {
	traceId := ctx.Value(consts.TraceId).(string)
	prefix := fmt.Sprintf("INFO  %s ", traceId)
	_ = log.Output(2, fmt.Sprintf(prefix+format+"\n", args...))
}

func WarnContextf(ctx context.Context, format string, args ...interface{}) {
	traceId := ctx.Value(consts.TraceId).(string)
	prefix := fmt.Sprintf("WARN  %s ", traceId)
	_ = log.Output(2, fmt.Sprintf(prefix+format+"\n", args...))
}

func ErrorContextf(ctx context.Context, format string, args ...interface{}) {
	traceId := ctx.Value(consts.TraceId).(string)
	prefix := fmt.Sprintf("ERROR %s ", traceId)
	_ = log.Output(2, fmt.Sprintf(prefix+format+"\n", args...))
}
