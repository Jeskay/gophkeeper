package control

type Controller interface {
	UploadFile(string) error
}
