package typesense

import (
	"mine/internal/models"
	"mine/internal/repositories"
	"mine/internal/settings"
)

type EventTypeSenseService interface {
	SearchText(dataInput models.Search) (interface{}, error)
}
type eventTypeSenseService struct {
	stt  *settings.AppSettings
	repo *repositories.Repositories
}

func NewTypeSenseService(
	appSettings *settings.AppSettings,
	repo *repositories.Repositories,
) EventTypeSenseService {
	return &eventTypeSenseService{
		stt:  appSettings,
		repo: repo,
	}
}
