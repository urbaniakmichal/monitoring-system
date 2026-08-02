package updates

type SystemUpdateInformation struct {
	HotFixID    string `json:"HotFixID"`
	Description string `json:"Description"`
	InstalledOn string `json:"InstalledOn"`
	InstalledBy string `json:"InstalledBy"`
}

type Updates interface {
	RetrieveSystemUpdates() ([]SystemUpdateInformation, error)
}

type SoftwareUpdates struct{}

var _ Updates = (*SoftwareUpdates)(nil)
