package storage

type Hero struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (h *Hero) IsValid() bool {
	if h.ID == "" {
		return false
	}
	if h.Name == "" {
		return false
	}
	return true
}

type Storager interface {
	Status() (string, error)
	GetHeroes() ([]Hero, error)
	GetHero(name string) (Hero, error)
	CreateHero(id, name string) error
	DeleteHero(id string) error
}

type ErrNothingToDelete struct {
	message string
}

func NewErrNothingToDelete(message string) *ErrNothingToDelete {
	return &ErrNothingToDelete{
		message: message,
	}
}
func (e *ErrNothingToDelete) Error() string {
	return e.message
}

type ErrHeroNotExist struct {
	message string
}

func NewErrHeroNotExist(message string) *ErrHeroNotExist {
	return &ErrHeroNotExist{
		message: message,
	}
}
func (e *ErrHeroNotExist) Error() string {
	return e.message
}

type ErrHeroRequestValidation struct {
	message string
}

func NewErrHeroRequestValidation(message string) *ErrHeroRequestValidation {
	return &ErrHeroRequestValidation{
		message: message,
	}
}
func (e *ErrHeroRequestValidation) Error() string {
	return e.message
}
