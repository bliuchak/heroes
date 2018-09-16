package storage

// ErrNothingToDelete custom error type for Hero handlers
// it tells that hero which requested to be deleted not exists
type ErrNothingToDelete struct {
	message string
}

// NewErrNothingToDelete returns pointer with error message to ErrNothingToDelete
func NewErrNothingToDelete(message string) *ErrNothingToDelete {
	return &ErrNothingToDelete{
		message: message,
	}
}

func (e *ErrNothingToDelete) Error() string {
	return e.message
}

// ErrHeroNotExist custom error for Hero handlers
// it tells that requested hero is not existing
type ErrHeroNotExist struct {
	message string
}

// NewErrHeroNotExist returns pointer with error message to ErrHeroNotExist
func NewErrHeroNotExist(message string) *ErrHeroNotExist {
	return &ErrHeroNotExist{
		message: message,
	}
}
func (e *ErrHeroNotExist) Error() string {
	return e.message
}
