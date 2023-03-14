package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/basebandit/bookshop/server/book/internal/controller/book"
	metadatagateway "github.com/basebandit/bookshop/server/book/internal/gateway/metadata/http"
	ratinggateway "github.com/basebandit/bookshop/server/book/internal/gateway/rating/http"
	httphandler "github.com/basebandit/bookshop/server/book/internal/handler/http"
	"github.com/basebandit/bookshop/server/pkg/discovery"
	"github.com/basebandit/bookshop/server/pkg/discovery/consul"
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
	metdataGateway := metadatagateway.New(registry)
	ratingGateway := ratinggateway.New(registry)
	svc := book.New(ratingGateway, metdataGateway)
	h := httphandler.New(svc)
	http.Handle("/book", http.HandlerFunc(h.GetBook))
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		panic(err)
	}
}
