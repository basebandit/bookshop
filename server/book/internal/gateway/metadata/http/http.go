package http

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/basebandit/bookshop/server/book/internal/gateway"
	"github.com/basebandit/bookshop/server/metadata/pkg/model"
)

// Gateway defines a book metadata HTTP gateway.
type Gateway struct {
	addr string
}

// New creates a new HTTP gateway for a book metadata service.
func New(addr string) *Gateway {
	return &Gateway{addr}
}

func (g *Gateway) Get(ctx context.Context, id string) (*model.Metadata, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%v/metadata", g.addr), nil)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	values := req.URL.Query()
	values.Add("id", id)
	req.URL.RawQuery = values.Encode()
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return nil, gateway.ErrNotFound
	} else if resp.StatusCode/100 != 2 {
		return nil, fmt.Errorf("non-2xx response: %v", resp)
	}

	var v *model.Metadata
	if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
		return nil, err
	}
	return v, nil
}