package dto

import (
	proto "gophkeeper/api/protos"
)

type File struct {
	Id     int64
	Name   string
	Status FileStatus
	Size   int
}

func (fs FileStatus) Proto() proto.FileInfo_Status {
	var status proto.FileStatus
	switch fs {
	case Available:
		status = proto.FileStatus_AVAILABLE
	case Reserved:
		status = proto.FileStatus_RESTRICTED
	default:
		status = proto.FileStatus_PROCESSING
	}
	return proto.FileInfo_Status{Status: status}
}

type FileStatus string

const (
	Available  = "AVAILABLE"
	Processing = "PROCESSING"
	Reserved   = "RESERVED"
)
