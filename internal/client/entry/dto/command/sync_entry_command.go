package command

import (
	"github.com/anoriar/gophkeeper/internal/client/entry/enum"
	validation "github.com/anoriar/gophkeeper/internal/client/shared/dto"
)

type SyncEntryCommand struct {
	EntryType enum.EntryType
}

func (s SyncEntryCommand) Validate() validation.ValidationErrors {
	return nil
}
