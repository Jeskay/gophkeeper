package abstraction

type FileData struct {
	Id     int64
	Name   string
	Size   string
	Status string
}

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
