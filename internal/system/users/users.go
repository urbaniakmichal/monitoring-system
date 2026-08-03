package users

type UserInformation struct {
	Username string `json:"Username"`
	Terminal string `json:"Terminal"`
	Host     string `json:"Host"`
}

type User interface {
	RetrieveUsersInfo() ([]UserInformation, error)
}

type SystemUsers struct{}

var _ User = (*SystemUsers)(nil)