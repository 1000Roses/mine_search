package repositories

import (
	control_api_tb "mine/internal/repositories/control_api_tb"
	sample_tb "mine/internal/repositories/sample"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Repositories struct {
	Sample     sample_tb.SampleRepo
	ControlApi control_api_tb.ControlApiRepo
}

func NewRepositories(
	db *gorm.DB,
	log *zap.Logger,
) *Repositories {
	return &Repositories{
		Sample:     sample_tb.NewSampleRepo(db, log),
		ControlApi: control_api_tb.NewControlApiRepo(db, log),
	}
}
