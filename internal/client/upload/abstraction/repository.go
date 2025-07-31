package abstraction

type uploadRepository struct {
	filePath string
	fileData []byte
}

func NewRepository() *uploadRepository {
	return &uploadRepository{}
}

func (r *uploadRepository) SetFilePath(value string) {
	r.filePath = value
}

func (r *uploadRepository) GetFilePath() string {
	return r.filePath
}

func (r *uploadRepository) SetFileData(value []byte) {
	r.fileData = value
}

func (r *uploadRepository) GetFileData() []byte {
	return r.fileData
}
