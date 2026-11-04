package main

type FirstResp struct {
	Response string `json:"response"`
}

type ApiKey struct {
	ID   string `json:"id"`
	Key  string `json:"key"`
	Name string `json:"name"`
}
