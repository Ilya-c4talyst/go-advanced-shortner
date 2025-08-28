package model

// Model Request
type Request struct {
	URL string `json:"url" validate:"required,url"`
}

// Model Response
type Response struct {
	Result string `json:"result"`
}

// Model for URL storage in JSON format
type URLRecord struct {
	UUID        string `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}
