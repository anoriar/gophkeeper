package command

import validation "github.com/anoriar/gophkeeper/internal/client/shared/dto"

type SyncEntryCommand struct {
}

func (s SyncEntryCommand) Validate() validation.ValidationErrors {
	return nil
}
