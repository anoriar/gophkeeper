package command

import (
	"github.com/anoriar/gophkeeper/internal/client/entry/enum"
	validation "github.com/anoriar/gophkeeper/internal/client/shared/dto"
)

type DeleteEntryCommand struct {
	Id        string
	EntryType enum.EntryType
}

func (d DeleteEntryCommand) Validate() validation.ValidationErrors {
	//TODO implement me
	panic("implement me")
}
