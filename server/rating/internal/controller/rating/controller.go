package rating

import (
	"context"
	"errors"

	"github.com/basebandit/bookshop/server/rating/internal/repository"
	"github.com/basebandit/bookshop/server/rating/pkg/model"
)

// ErrNotFound is returned when no ratings are found for a record.
var ErrNotFound = errors.New("ratings not found for a record")

type ratingRepository interface {
	Get(ctx context.Context, recordID model.RecordID, recordType model.RecordType) ([]model.Rating, error)
	Add(ctx context.Context, recordID model.RecordID, recordType model.RecordType, rating *model.Rating) error
}

// Controller defines a rating service controller.
type Controller struct {
	repo ratingRepository
}

// New creates a rating service controller
func New(repo ratingRepository) *Controller {
	return &Controller{repo}
}

// GetAggregateRating returns the aggregated rating for a record or ErrNotFound if there are no ratings for it.
func (c *Controller) GetAggregateRating(ctx context.Context, recordID model.RecordID, recordType model.RecordType) (float64, error) {
	ratings, err := c.repo.Get(ctx, recordID, recordType)
	if err != nil && err == repository.ErrNotFound {
		return 0, err
	} else if err != nil {
		return 0, err
	}
	var sum float64
	for _, r := range ratings {
		sum += float64(r.Value)
	}
	return sum / float64(len(ratings)), nil
}

// AddRating writes a rating for a given record.
func (c *Controller) AddRating(ctx context.Context, recordID model.RecordID, recordType model.RecordType, rating *model.Rating) error {
	return c.repo.Add(ctx, recordID, recordType, rating)
}
