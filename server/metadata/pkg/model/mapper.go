package model

import "github.com/basebandit/bookshop/server/gen"

// MetadataToProto converts a Metadata struct into a generated proto counterpart.
func MetadataToProto(m *Metadata) *gen.Metadata {
	return &gen.Metadata{
		Id:                m.ID,
		Title:             m.Title,
		Author:            m.Author,
		Description:       m.Description,
		PublishingYear:    m.PublishingYear,
		PublishingCompany: m.PublishingCompany,
	}
}

// MetadataFromProto converts a generated proto counterpart into a Metadata struct
func MetadataFromProto(m *gen.Metadata) *Metadata {
	return &Metadata{
		ID:                m.Id,
		Title:             m.Title,
		Author:            m.Author,
		Description:       m.Description,
		PublishingYear:    m.PublishingYear,
		PublishingCompany: m.PublishingCompany,
	}
}
