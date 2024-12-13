package sample

import (
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type (
	SampleRepo interface{}
	sampleRepo struct {
		db  *gorm.DB
		log *zap.Logger
	}
)

func NewSampleRepo(
	db *gorm.DB,
	log *zap.Logger,
) SampleRepo {
	return &sampleRepo{
		db:  db,
		log: log,
	}
}
