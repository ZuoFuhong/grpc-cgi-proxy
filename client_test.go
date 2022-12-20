package main

import (
	"context"
	"fmt"
	codec "github.com/ZuoFuhong/grpc-middleware/encoding/json"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"testing"
)

func Test_RpcInvoke(t *testing.T) {
	conn, err := grpc.Dial("127.0.0.1:1025", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatal(err)
	}
	req := map[string]string{"private_key": "0x01c4bda0939df07a31e3738c6c1e1d5905c9f229e6ffa1922557308a62efb23f"}
	rsp := make(map[string]string)
	if err := conn.Invoke(context.Background(), "/go_wallet_manage_svr.go_wallet_manage_svr/ImportWallet", req, &rsp, grpc.CallContentSubtype(codec.Name)); err != nil {
		t.Fatal(err)
	}
	fmt.Println(rsp)
}
