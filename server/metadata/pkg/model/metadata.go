package model

// Metadata defines the book metadata.
type Metadata struct {
	ID                string `json:"id"`
	Title             string `json:"title"`
	Description       string `json:"description"`
	Author            string `json:"author"`
	PublishingYear    string `json:"publishingYear"`
	PublishingCompany string `json:"publishingCompany"`
}
