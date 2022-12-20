package main

import (
	"github.com/ZuoFuhong/grpc-cgi-proxy/server"
	"log"
)

func main() {
	log.SetFlags(log.Lshortfile | log.Ltime | log.Ldate)
	if err := server.NewServer().Serve(); err != nil {
		log.Fatal(err)
	}
}
