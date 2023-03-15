package grpc

import (
	"context"
	"errors"

	"github.com/basebandit/bookshop/server/gen"
	"github.com/basebandit/bookshop/server/metadata/internal/controller/metadata"
	"github.com/basebandit/bookshop/server/metadata/pkg/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Handler defines a book metadata gRPC handler.
type Handler struct {
	gen.UnimplementedMetadataServiceServer
	ctrl *metadata.Controller
}

// New creates a new book metadata gRPC handler
func New(ctrl *metadata.Controller) *Handler {
	return &Handler{ctrl: ctrl}
}

// GetMetadata returns book metadata.
func (h *Handler) GetMetadata(ctx context.Context, req *gen.GetMetadataRequest) (*gen.GetMetadataResponse, error) {
	if req == nil || req.Id == "" {
		return nil, status.Errorf(codes.InvalidArgument, "nil req or empty id")
	}
	book, err := h.ctrl.Get(ctx, req.Id)
	if err != nil && errors.Is(err, metadata.ErrNotFound) {
		return nil, status.Errorf(codes.NotFound, err.Error())
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &gen.GetMetadataResponse{Metadata: model.MetadataToProto(book)}, nil
}
