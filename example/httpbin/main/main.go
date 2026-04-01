package main

import (
	"context"
	"fmt"
	"github.com/mono83/transport/example/httpbin"
	"github.com/mono83/transport/http/log"
	"github.com/mono83/transport/http/native"
)

func main() {
	t := native.NewWithLog(log.Stdout)
	hb := httpbin.Client{Transport: t}

	ctx := context.Background()
	//ctx, _ = context.WithTimeout(ctx, time.Millisecond)
	fmt.Println(httpbin.PostCall{Name: "foo", Value: 99}.Execute(ctx, hb))
}
