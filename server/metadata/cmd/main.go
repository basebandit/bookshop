package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/basebandit/bookshop/server/gen"
	"github.com/basebandit/bookshop/server/metadata/internal/controller/metadata"
	grpchandler "github.com/basebandit/bookshop/server/metadata/internal/handler/grpc"
	"github.com/basebandit/bookshop/server/metadata/internal/repository/memory"
	"github.com/basebandit/bookshop/server/pkg/discovery"
	"github.com/basebandit/bookshop/server/pkg/discovery/consul"
	"google.golang.org/grpc"
)

const serviceName = "metadata"

func main() {
	var port int
	flag.IntVar(&port, "port", 8081, "API handler port")
	flag.Parse()
	log.Printf("starting the metadata service on port %d\n", port)
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
	repo := memory.New()
	ctrl := metadata.New(repo)
	h := grpchandler.New(ctrl)
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	srv := grpc.NewServer()
	gen.RegisterMetadataServiceServer(srv, h)
	if err := srv.Serve(lis); err != nil {
		panic(err)
	}
}
