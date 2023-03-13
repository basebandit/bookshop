package book

import (
	"context"
	"errors"

	"github.com/basebandit/bookshop/server/book/internal/gateway"
	"github.com/basebandit/bookshop/server/book/pkg/model"
	metadatamodel "github.com/basebandit/bookshop/server/metadata/pkg/model"
	ratingmodel "github.com/basebandit/bookshop/server/rating/pkg/model"
)

// ErrNotFound is returned when the book metadata is not found.
var ErrNotFound = errors.New("book metadata not found")

type ratingGateway interface {
	GetAggregatedRating(ctx context.Context, recordID ratingmodel.RecordID, recordType ratingmodel.RecordType) (float64, error)
	AddRating(ctx context.Context, recordID ratingmodel.RecordID, recordType ratingmodel.RecordType, rating *ratingmodel.Rating) error
}

type metadataGateway interface {
	Get(ctx context.Context, id string) (*metadatamodel.Metadata, error)
}

// Controller defines a book service controller
type Controller struct {
	ratingGateway   ratingGateway
	metadataGateway metadataGateway
}

// New creates a new book service controller
func New(ratingGateway ratingGateway, metadataGateway metadataGateway) *Controller {
	return &Controller{ratingGateway, metadataGateway}
}

// Get returns the book details including the aggregated rating and metadata.
func (c *Controller) Get(ctx context.Context, id string) (*model.Book, error) {
	metadata, err := c.metadataGateway.Get(ctx, id)
	if err != nil && errors.Is(err, gateway.ErrNotFound) {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}
	book := &model.Book{Metadata: *metadata}
	rating, err := c.ratingGateway.GetAggregatedRating(ctx, ratingmodel.RecordID(id), ratingmodel.RecordTypeBook)
	if err != nil && !errors.Is(err, gateway.ErrNotFound) {
		return nil, err
	} else if err != nil {
		return nil, err
	} else {
		book.Rating = &rating
	}
	return book, nil
}
