package dto

import "github.com/bastianvv/vio/internal/domain"

type Library struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

func NewLibrary(l *domain.Library) *Library {
	return &Library{
		ID:   l.ID,
		Name: l.Name,
		Type: string(l.Type), // enum: movies | series | anime | others
	}
}
