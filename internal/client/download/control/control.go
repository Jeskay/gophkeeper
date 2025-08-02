package control

import (
	"gophkeeper/internal/client/download/abstraction"
)

type Controller interface {
	DownloadFile(id int64) error
	DownloadFileList() ([]*abstraction.FileData, error)
}
