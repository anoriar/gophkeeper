package command

import (
	"github.com/anoriar/gophkeeper/internal/client/entry/enum"
	validation "github.com/anoriar/gophkeeper/internal/client/shared/dto"
)

type DetailEntryCommand struct {
	Id        string
	EntryType enum.EntryType
}

func (d DetailEntryCommand) Validate() validation.ValidationErrors {
	//TODO implement me
	panic("implement me")
}
