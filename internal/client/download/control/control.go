package control

import (
	"gophkeeper/internal/client/download/abstraction"
)

type Controller interface {
	DownloadFile(fileName string) error
	DownloadFileList() ([]*abstraction.FileData, error)
}
