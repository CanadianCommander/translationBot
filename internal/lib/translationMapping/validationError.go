package translationMapping

//==========================================================================
// Public
//==========================================================================

type ValidationError struct {
	Msg string
}

// NewValidationError creates a csv validation error with the given message
func NewValidationError(msg string) ValidationError {
	return ValidationError{
		Msg: msg,
	}
}

func (c ValidationError) Error() string {
	return c.Msg
}
