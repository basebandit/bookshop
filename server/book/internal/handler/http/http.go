package http

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/basebandit/bookshop/server/book/internal/controller/book"
)

// Handler defines a book service handler.
type Handler struct {
	ctrl *book.Controller
}

// New creates a new book HTTP handler.
func New(ctrl *book.Controller) *Handler {
	return &Handler{ctrl}
}

// GetBook handles GET /book requests.
func (h *Handler) GetBook(w http.ResponseWriter, req *http.Request) {
	id := req.FormValue("id")
	bookDetails, err := h.ctrl.Get(req.Context(), id)
	if err != nil && errors.Is(err, book.ErrNotFound) {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		log.Printf("repository get error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(bookDetails); err != nil {
		log.Printf("response encode error: %v\n", err)
	}
}
