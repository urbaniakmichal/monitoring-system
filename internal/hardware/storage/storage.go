package storage

type StorageInformation struct {
	Device      string  `json:"Device"`
	Path        string  `json:"Path"`
	TotalMB     uint64  `json:"TotalMB"`
	FreeMB      uint64  `json:"FreeMB"`
	UsedPercent float64 `json:"UsedPercent"`
}

type Storage interface {
	RetrieveStorageInformation() ([]StorageInformation, error)
}

type HardwareStorage struct{}

var _ Storage = (*HardwareStorage)(nil)