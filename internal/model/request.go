package model

// Model Request
type Request struct {
	URL string `json:"url" validate:"required,url"`
}

// Model Response
type Response struct {
	Result string `json:"result"`
}
