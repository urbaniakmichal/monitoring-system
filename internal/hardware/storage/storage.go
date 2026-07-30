package storage

type StorageInformation struct {
	Path        string  `json:"Path"`
	TotalMB     uint64  `json:"TotalMB"`
	FreeMB      uint64  `json:"FreeMB"`
	UsedPercent float64 `json:"UsedPercent"`
}