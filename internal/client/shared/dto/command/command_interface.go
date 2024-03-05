package command

import validation "github.com/anoriar/gophkeeper/internal/client/shared/dto"

type CommandInterface interface {
	Validate() validation.ValidationErrors
}
