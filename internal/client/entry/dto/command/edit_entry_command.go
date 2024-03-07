package command

import (
	"encoding/json"
	"fmt"

	"github.com/anoriar/gophkeeper/internal/client/entry/dto"

	"github.com/anoriar/gophkeeper/internal/client/entry/enum"
	validation "github.com/anoriar/gophkeeper/internal/client/shared/dto"
)

type EditEntryCommand struct {
	Id        string
	EntryType enum.EntryType
	Data      interface{}
	Meta      json.RawMessage
}

func (command *EditEntryCommand) Validate() validation.ValidationErrors {
	var validationErrors validation.ValidationErrors

	if command.Id == "" {
		validationErrors = append(validationErrors, fmt.Errorf("id required"))
	}

	switch command.EntryType {
	case enum.Login:
		data, ok := command.Data.(dto.LoginData)
		if !ok {
			validationErrors = append(validationErrors, fmt.Errorf("data not compatible with any format"))
		}
		validationErrors = append(validationErrors, data.Validate()...)
	case enum.Card:
		data, ok := command.Data.(dto.CardData)
		if !ok {
			validationErrors = append(validationErrors, fmt.Errorf("data not compatible with any format"))
		}
		validationErrors = append(validationErrors, data.Validate()...)
	default:
		validationErrors = append(validationErrors, fmt.Errorf("data not compatible with any format"))
	}
	return nil
}
