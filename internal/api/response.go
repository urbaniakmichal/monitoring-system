package api

type Link struct {
	Href   string `json:"href"`
	Method string `json:"method"`
}

type APIResponse struct {
	Data  any             `json:"data,omitempty"`
	Links map[string]Link `json:"_links"`
}

type ErrorResponse struct {
	Error string          `json:"error"`
	Links map[string]Link `json:"_links"`
}
