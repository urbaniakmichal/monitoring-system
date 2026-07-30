package peripherals

type DisplayInformation struct {
	Name       string `json:"Name"`
	Resolution string `json:"Resolution"`
	IsPrimary  bool   `json:"IsPrimary"`
}

type USBDeviceInformation struct {
	Name         string `json:"Name"`
	Manufacturer string `json:"Manufacturer"`
	DeviceID     string `json:"DeviceID"`
}

type PeripheralsInformation struct {
	Displays []DisplayInformation   `json:"Displays"`
	USB      []USBDeviceInformation `json:"USB"`
}