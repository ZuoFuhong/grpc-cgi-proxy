package server

import (
	"fmt"
	"github.com/ZuoFuhong/grpc-cgi-proxy/pkg/config"
	"github.com/ZuoFuhong/grpc-cgi-proxy/pkg/log"
	"github.com/ZuoFuhong/grpc-naming-monica/registry"
	"net/http"
)

type Server struct {
	name   string
	addr   string
	router *Router
}

func NewServer() *Server {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic("load config fail: " + err.Error())
	}
	config.SetGlobalConfig(cfg)

	// 服务注册
	if err := registry.NewRegistry(&registry.Config{
		Token:       cfg.Monica.Token,
		Namespace:   cfg.Monica.Namespace,
		ServiceName: cfg.Monica.ServiceName,
		IP:          cfg.Server.Addr,
		Port:        cfg.Server.Port,
	}).Register(); err != nil {
		log.Fatal(err)
	}

	return &Server{
		name:   cfg.Server.Name,
		addr:   fmt.Sprintf("%s:%d", cfg.Server.Addr, cfg.Server.Port),
		router: NewRouter(),
	}
}

func (s *Server) Serve() error {
	log.Debugf("%s runs on http://%s", s.name, s.addr)
	return http.ListenAndServe(s.addr, s.router)
}
