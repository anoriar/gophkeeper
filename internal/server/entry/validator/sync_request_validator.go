package validator

import (
	"fmt"

	"github.com/anoriar/gophkeeper/internal/server/entry/dto/request/sync"
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
		if reqItem.OriginalId == "" {
			validationErrors = append(validationErrors, fmt.Errorf("item %d: originalId required", itemIndex))
		}
		if len(reqItem.Data) == 0 {
			validationErrors = append(validationErrors, fmt.Errorf("item %d: data required", itemIndex))
		}
	}
	return validationErrors
}
