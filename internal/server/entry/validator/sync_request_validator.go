package validator

import (
	"fmt"

	"github.com/anoriar/gophkeeper/internal/server/entry/dto"
	"github.com/anoriar/gophkeeper/internal/server/entry/dto/request/sync"
	"github.com/anoriar/gophkeeper/internal/server/entry/enum"
	validation "github.com/anoriar/gophkeeper/internal/server/shared/dto"
)

type SyncRequestValidator struct {
}

func NewSyncRequestValidator() *SyncRequestValidator {
	return &SyncRequestValidator{}
}

func (v *SyncRequestValidator) ValidateSyncRequest(request sync.SyncRequest) validation.ValidationErrors {

	var validationErrors validation.ValidationErrors
	for itemIndex, reqItem := range request.Items {
		if reqItem.Id == "" {
			validationErrors = append(validationErrors, fmt.Errorf("item %d: id required", itemIndex))
		}
		if !enum.IsEntryType(reqItem.EntryType) {
			validationErrors = append(validationErrors, fmt.Errorf("item %d: not valid entry type", itemIndex))
		}

		switch reqItem.EntryType {
		case enum.Login:
			data, ok := reqItem.Data.(*dto.LoginData)
			if !ok {
				validationErrors = append(validationErrors, fmt.Errorf("item %d: data not compatible with any format", itemIndex))
			}
			validationErrors = append(validationErrors, data.Validate()...)
		case enum.Card:
			data, ok := reqItem.Data.(*dto.CardData)
			if !ok {
				validationErrors = append(validationErrors, fmt.Errorf("item %d: data not compatible with any format", itemIndex))
			}
			validationErrors = append(validationErrors, data.Validate()...)
		default:
			validationErrors = append(validationErrors, fmt.Errorf("item %d: data not compatible with any format", itemIndex))
		}

	}
	return validationErrors
}
