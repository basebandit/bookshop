package grpc

import (
	"context"
	"errors"

	"github.com/basebandit/bookshop/server/book/internal/controller/book"
	"github.com/basebandit/bookshop/server/gen"
	"github.com/basebandit/bookshop/server/metadata/pkg/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Handler defines a book gRPC handler.
type Handler struct {
	gen.UnimplementedBookServer
	ctrl *book.Controller
}

// New creates a new book gRPC handler
func New(ctrl *book.Controller) *Handler {
	return &Handler{ctrl: ctrl}
}

// GetBook returns book details using the given id
func (h *Handler) GetBook(ctx context.Context, req *gen.GetBookRequest) (*gen.GetBookResponse, error) {
	if req == nil || req.BookId == "" {
		return nil, status.Errorf(codes.InvalidArgument, "nil req or empty id")
	}
	bookDetails, err := h.ctrl.Get(ctx, req.BookId)
	if err != nil && errors.Is(err, book.ErrNotFound) {
		return nil, status.Errorf(codes.NotFound, err.Error())
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &gen.GetBookResponse{
		Book: &gen.BookDetails{
			Metadata: model.MetadataToProto(&bookDetails.Metadata),
			Rating:   *bookDetails.Rating,
		},
	}, nil
}
