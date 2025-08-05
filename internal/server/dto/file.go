package dto

import (
	proto "gophkeeper/api/protos"
)

type ChunkData struct {
	Data []byte
	Err  error
}

type FileStatus string

const (
	Available  = "AVAILABLE"
	Processing = "PROCESSING"
	Reserved   = "RESERVED"
)

type File struct {
	Id     int64
	Name   string
	Status FileStatus
	Size   int
}

func (f *File) Proto() *proto.FileInfo {
	return &proto.FileInfo{
		Name:           f.Name,
		OptionalStatus: f.Status.Proto(),
		OptionalId:     &proto.FileInfo_Id{Id: f.Id},
		Size:           uint32(f.Size),
	}
}

func (fs *FileStatus) Proto() *proto.FileInfo_Status {
	var status proto.FileStatus
	switch *fs {
	case Available:
		status = proto.FileStatus_AVAILABLE
	case Reserved:
		status = proto.FileStatus_RESTRICTED
	default:
		status = proto.FileStatus_PROCESSING
	}
	return &proto.FileInfo_Status{Status: status}
}
