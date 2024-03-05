package command

import (
	"github.com/anoriar/gophkeeper/internal/client/entry/enum"
	validation "github.com/anoriar/gophkeeper/internal/client/shared/dto"
)

type ListEntryCommand struct {
	EntryType enum.EntryType
}

func (l ListEntryCommand) Validate() validation.ValidationErrors {
	//TODO implement me
	panic("implement me")
}
