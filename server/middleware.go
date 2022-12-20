package server

import (
	"context"
	"github.com/ZuoFuhong/grpc-cgi-proxy/consts"
	"github.com/ZuoFuhong/grpc-cgi-proxy/pkg/log"
	"github.com/google/uuid"
	"net/http"
)

type Middleware struct {
}

func NewMiddleware() *Middleware {
	return &Middleware{}
}

func (m *Middleware) RequestMetricHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		traceId := uuid.New().String()
		r = r.WithContext(context.WithValue(r.Context(), consts.TraceId, traceId))
		log.Debugf("%s %s %q", traceId, r.Method, r.URL.String())
		w.Header().Add("request_id", traceId)
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
