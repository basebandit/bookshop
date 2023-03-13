package model

import "github.com/basebandit/bookshop/server/metadata/pkg/model"

type Book struct {
	Rating   *float64       `json:"rating,omitempty"`
	Metadata model.Metadata `json:"metadata"`
}
