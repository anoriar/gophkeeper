package command

import (
	"fmt"

	"github.com/anoriar/gophkeeper/internal/client/entry/enum"
	validation "github.com/anoriar/gophkeeper/internal/client/shared/dto"
)

type DetailEntryCommand struct {
	Id        string
	EntryType enum.EntryType
}

func (command *DetailEntryCommand) Validate() validation.ValidationErrors {
	var validationErrors validation.ValidationErrors
	if command.Id == "" {
		validationErrors = append(validationErrors, fmt.Errorf("id required"))
	}
	return nil
}
