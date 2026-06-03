package request

type LlmConfigReq struct {
	Provider    string  `json:"provider"`
	APIURL      string  `json:"api_url"`
	APIKey      string  `json:"api_key"`
	Model       string  `json:"model"`
	MaxTokens   int     `json:"max_tokens"`
	Temperature float64 `json:"temperature"`
}
