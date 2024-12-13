package sample

import (
	"mine/internal/repositories"
	"mine/internal/settings"
)

type EventSampleService interface {
	DoFuncSample(phone, mail string) (interface{}, error)
}
type eventSampleService struct {
	stt  *settings.AppSettings
	repo *repositories.Repositories
}

func NewSampleService(
	appSettings *settings.AppSettings,
	repo *repositories.Repositories,
) EventSampleService {
	return &eventSampleService{
		stt:  appSettings,
		repo: repo,
	}
}
