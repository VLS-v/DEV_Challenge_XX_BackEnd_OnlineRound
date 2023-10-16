package models

type Cell struct {
	Value  string `json:"value"`
	Result string `json:"result"`
}

type Sheet map[string]*Cell

type SetCell struct {
	Value string `json:"value"`
}
