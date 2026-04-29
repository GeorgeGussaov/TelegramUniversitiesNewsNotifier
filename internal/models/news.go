package models

type News struct {
	Source      string   `json:"Source"`
	Title       string   `json:"Title"`
	Date        string   `json:"Date"`
	Link        string   `json:"Link"`
	Text        string   `json:"Text"`
	ImagesLinks []string `json:"ImagesLinks"`
}