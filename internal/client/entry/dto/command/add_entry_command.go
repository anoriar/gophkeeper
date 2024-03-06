package command

import (
	"encoding/json"

	"github.com/anoriar/gophkeeper/internal/client/entry/enum"
	validation "github.com/anoriar/gophkeeper/internal/client/shared/dto"
)

type AddEntryCommand struct {
	EntryType enum.EntryType
	Data      interface{}
	Meta      json.RawMessage
}

func (command *AddEntryCommand) Validate() validation.ValidationErrors {
	return nil
}
