package grpc

import (
	"context"

	"github.com/basebandit/bookshop/server/gen"
	"github.com/basebandit/bookshop/server/internal/grpcutil"
	"github.com/basebandit/bookshop/server/metadata/pkg/model"
	"github.com/basebandit/bookshop/server/pkg/discovery"
)

// Gateway defines a book metadata gRPC gateway.
type Gateway struct {
	registry discovery.Registry
}

// New creates a new gRPC gateway for a book metadata service
func New(registry discovery.Registry) *Gateway {
	return &Gateway{registry}
}

// Get returns book metadata using the given id
func (g *Gateway) Get(ctx context.Context, id string) (*model.Metadata, error) {
	conn, err := grpcutil.ServiceConnection(ctx, "metadata", g.registry)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	client := gen.NewMetadataServiceClient(conn)
	resp, err := client.GetMetadata(ctx, &gen.GetMetadataRequest{
		Id: id,
	})
	if err != nil {
		return nil, err
	}
	return model.MetadataFromProto(resp.Metadata), nil
}
