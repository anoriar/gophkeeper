package command

import (
	"encoding/json"
	"github.com/anoriar/gophkeeper/internal/client/entry/enum"
	validation "github.com/anoriar/gophkeeper/internal/client/shared/dto"
)

type EditEntryCommand struct {
	Id        string
	EntryType enum.EntryType
	Data      interface{}
	Meta      json.RawMessage
}

func (e EditEntryCommand) Validate() validation.ValidationErrors {
	//TODO implement me
	panic("implement me")
}
