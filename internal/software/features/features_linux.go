//go:build linux

package features

type OptionalFeatureInformation struct {
	FeatureName string `json:"FeatureName"`
	State       string `json:"State"`
}

func RetrieveOptionalFeatures() ([]OptionalFeatureInformation, error) {
	return []OptionalFeatureInformation{}, nil
}