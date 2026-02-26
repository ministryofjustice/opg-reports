package apimodels

// Config contains required values for DB and others to generate a response
type Args struct {
	DB      string `json:"db"`
	Driver  string `json:"driver"`
	Params  string `json:"params"`
	Version string `json:"version"`
	SHA     string `json:"sha"`
}
