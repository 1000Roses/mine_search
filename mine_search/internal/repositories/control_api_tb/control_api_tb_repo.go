package control_api_tb

import (
	"mine/internal/models"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ControlApiRepo interface {
	GetInfoFrApiName(string) (*models.ControlApiTb, error)
}
type controlApiRepo struct {
	db  *gorm.DB
	log *zap.Logger
}

func NewControlApiRepo(
	db *gorm.DB,
	log *zap.Logger,
) ControlApiRepo {
	return &controlApiRepo{
		db:  db,
		log: log,
	}
}
