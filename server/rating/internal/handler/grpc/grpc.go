package grpc

import (
	"context"
	"errors"

	"github.com/basebandit/bookshop/server/gen"
	"github.com/basebandit/bookshop/server/rating/internal/controller/rating"
	"github.com/basebandit/bookshop/server/rating/pkg/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Handler defines a gRPC rating API handler.
type Handler struct {
	gen.UnimplementedRatingServer
	ctrl *rating.Controller
}

// New creates a new book rating gRPC handler.
func New(ctrl *rating.Controller) *Handler {
	return &Handler{ctrl: ctrl}
}

// GetAggregatedRating returns the aggregated rating for a record.
func (h *Handler) GetAggregatedRating(ctx context.Context, req *gen.GetAggregatedRatingRequest) (*gen.GetAggregatedRatingResponse, error) {
	if req == nil || req.RecordId == "" || req.RecordType == "" {
		return nil, status.Errorf(codes.InvalidArgument, "nil req or empty id")
	}
	v, err := h.ctrl.GetAggregateRating(ctx, model.RecordID(req.RecordId), model.RecordType(req.RecordType))
	if err != nil && errors.Is(err, rating.ErrNotFound) {
		return nil, status.Errorf(codes.NotFound, err.Error())
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &gen.GetAggregatedRatingResponse{RatingValue: v}, nil
}

// AddRating writes a rating for a given book record.
func (h *Handler) AddRating(ctx context.Context, req *gen.AddRatingRequest) (*gen.AddRatingResponse, error) {
	if req == nil || req.RecordId == "" || req.UserId == "" {
		return nil, status.Errorf(codes.InvalidArgument, "nil req or empty user id or record id")
	}
	if err := h.ctrl.AddRating(ctx, model.RecordID(req.RecordId), model.RecordType(req.RecordType), &model.Rating{
		UserID: model.UserID(req.UserId), Value: model.RatingValue(req.RatingValue),
	}); err != nil {
		return nil, err
	}
	return &gen.AddRatingResponse{}, nil
}
