package typesense

import (
	"errors"
	"mine/internal"
	"mine/internal/models"
	utilsCall "mine/internal/utils_call"

	"go.uber.org/zap"
)

func (a *eventTypeSenseService) SearchText(dataInput models.Search) (interface{}, error) {
	document, query := "", ""
	result, err := utilsCall.TypeSenseSearchText(a.stt, a.repo, document, query)
	if err != nil {
		internal.Log.Error("SearchText -> TypeSenseSearchText", zap.Any("document", document), zap.Any("query", query), zap.Error(err))
		return nil, errors.New(internal.SysStatus.SystemError.Msg)
	}
	return result, nil
}
