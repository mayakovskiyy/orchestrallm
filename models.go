package main

type FirstResp struct {
	Response string `json:"response"`
}

type ApiKey struct {
	ID   string `json:"id"`
	Key  string `json:"key"`
	Name string `json:"name"`
}

type ModelSignal struct {
	ID     string `json:"id"`
	Model  string `json:"model"`
	Role   string `json:"role"`
	Prompt string `json:"prompt"`
}
