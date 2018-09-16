package storage

// Storager general storage interface
type Storager interface {
	Status() (string, error)
	GetHeroes() ([]Hero, error)
	GetHero(name string) (Hero, error)
	CreateHero(id, name string) error
	DeleteHero(id string) error
}

// Hero contains hero data
type Hero struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// IsValid validates hero structure
func (h *Hero) IsValid() bool {
	if h.ID == "" {
		return false
	}
	if h.Name == "" {
		return false
	}
	return true
}
