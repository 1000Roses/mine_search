package control_api_tb

import (
	"errors"
	"fmt"
	"mine/internal/models"

	"gorm.io/gorm"
)

func (repo *controlApiRepo) GetInfoFrApiName(apiname string) (*models.ControlApiTb, error) {
	result := &models.ControlApiTb{}
	where := fmt.Sprintf("%s = ?", result.ColumnApiName())
	err := repo.db.Debug().Table(result.TableName()).Where(where, apiname).First(&result).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return result, nil
}
