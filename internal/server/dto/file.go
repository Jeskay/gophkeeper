package dto

type File struct {
	Id     int64
	Name   string
	Status FileStatus
}

type FileStatus string

const (
	Available  = "AVAILABLE"
	Processing = "PROCESSING"
	Reserved   = "RESERVED"
)
