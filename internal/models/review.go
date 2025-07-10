package models

type Review struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Photo       string `json:"photo"`
	Description string `json:"description"`
	Rating      int    `json:"rating"`
}
