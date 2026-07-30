package io

type NetworkIO struct {
	BytesSent uint64 `json:"BytesSent"`
	BytesRecv uint64 `json:"BytesRecv"`
}

type DiskIO struct {
	ReadBytes  uint64 `json:"ReadBytes"`
	WriteBytes uint64 `json:"WriteBytes"`
}

type IOStatistics struct {
	Network map[string]NetworkIO `json:"Network"`
	Disk    map[string]DiskIO    `json:"Disk"`
}