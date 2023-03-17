package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/basebandit/bookshop/server/book/internal/controller/book"
	metadatagateway "github.com/basebandit/bookshop/server/book/internal/gateway/metadata/http"
	ratinggateway "github.com/basebandit/bookshop/server/book/internal/gateway/rating/http"
	grpchandler "github.com/basebandit/bookshop/server/book/internal/handler/grpc"
	"github.com/basebandit/bookshop/server/gen"
	"github.com/basebandit/bookshop/server/pkg/discovery"
	"github.com/basebandit/bookshop/server/pkg/discovery/consul" // Copyright  Canonical Ltd.
	"google.golang.org/grpc"
	// Licensed under the AGPLv3, see LICENCE file for details.
)

const serviceName = "book"

func main() {
	var port int
	flag.IntVar(&port, "port", 8083, "API handler port")
	flag.Parse()
	log.Printf("starting the book service on port %d\n", port)
	registry, err := consul.NewRegistry("localhost:8500")
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	instanceID := discovery.GenerateInstanceID(serviceName)
	if err := registry.Register(ctx, instanceID, serviceName, fmt.Sprintf("localhost:%d", port)); err != nil {
		panic(err)
	}
	go func() {
		for {
			if err := registry.ReportHealthyState(instanceID, serviceName); err != nil {
				log.Printf("failed to report healthy state: %v\n", err.Error())
			}
			time.Sleep(1 * time.Second)
		}
	}()
	defer func() {
		if err := registry.Deregister(ctx, instanceID, serviceName); err != nil {
			panic(err)
		}
	}()
	metadataGateway := metadatagateway.New(registry)
	ratingGateway := ratinggateway.New(registry)
	ctrl := book.New(ratingGateway, metadataGateway)
	h := grpchandler.New(ctrl)
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v\n", err)
	}
	srv := grpc.NewServer()
	gen.RegisterBookServer(srv, h)
	if err := srv.Serve(lis); err != nil {
		panic(err)
	}
}
