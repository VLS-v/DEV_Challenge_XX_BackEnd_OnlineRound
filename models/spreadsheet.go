package models

type Cell struct {
	Value string
	//Result string `json:"result"`
}

type Sheet map[string]Cell

type SavesData map[string]Sheet

type CellResponse struct {
	Value  string `json:"value"`
	Result string `json:"result"`
}

type SheetResponse map[string]CellResponse
