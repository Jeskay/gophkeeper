package file

type Service interface {
	SaveFile(name string, data []byte) error
	ReadFile(name string) ([]byte, error)
}
