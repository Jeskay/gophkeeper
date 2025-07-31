package abstraction

type menuRepository struct {
	availability map[string]bool
	token        string
}

func NewRepository() *menuRepository {
	return &menuRepository{}
}
func (r *menuRepository) SaveAvailability(value map[string]bool) {
	r.availability = value
}
func (r *menuRepository) GetAvailability() map[string]bool {
	return r.availability
}
func (r *menuRepository) SaveToken(value string) {
	r.token = value
}
func (r *menuRepository) GetToken() string {
	return r.token
}
