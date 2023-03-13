package memory

import (
	"context"
	"sync"

	"github.com/basebandit/bookshop/server/metadata/internal/repository"
	"github.com/basebandit/bookshop/server/metadata/pkg/model"
)

// Repository defines a memory book metadata repository
type Repository struct {
	sync.RWMutex
	data map[string]*model.Metadata
}

// New creates a new memory repository
func New() *Repository {
	return &Repository{data: make(map[string]*model.Metadata)}
}

// Get retrieves book metadata by book id.
func (r *Repository) Get(_ context.Context, id string) (*model.Metadata, error) {
	r.RLock()
	defer r.Unlock()

	m, ok := r.data[id]
	if !ok {
		return nil, repository.ErrNotFound
	}
	return m, nil
}

// Add adds book metadata for a given book id.
func (r *Repository) Add(_ context.Context, id string, metadata *model.Metadata) error {
	r.Lock()
	defer r.Unlock()
	r.data[id] = metadata
	return nil
}
