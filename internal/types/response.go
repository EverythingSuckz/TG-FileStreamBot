package types

type RootResponse struct {
	Message string `json:"message"`
	Ok      bool   `json:"ok"`
	Uptime  string `json:"uptime"`
	Version string `json:"version"`
}
