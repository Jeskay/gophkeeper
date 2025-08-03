package abstraction

type downloadRepository struct {
	files []*FileData
}

func NewRepository() *downloadRepository {
	return &downloadRepository{}
}

func (r *downloadRepository) SetFiles(value []*FileData) {
	r.files = value
}

func (r *downloadRepository) GetFiles() []*FileData {
	return r.files
}
