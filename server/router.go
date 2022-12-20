package server

import (
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"net/http"
)

type Router struct {
	*mux.Router
	ch *alice.Chain
}

func NewRouter() *Router {
	mw := NewMiddleware()
	chain := alice.New(mw.RequestMetricHandler)
	r := &Router{
		Router: mux.NewRouter(),
		ch:     &chain,
	}
	r.registerHandler()
	return r
}

func (r *Router) registerHandler() {
	r.Handle("/ping", r.ch.ThenFunc(Ping)).Methods("GET")
	r.PathPrefix("/cgi-proxy").Handler(r.ch.Then(NewDispatcher()))
}

func Ping(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("ok"))
}
