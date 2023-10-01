package models

type Cell struct {
	Value  string
	Result string
}

type Sheet map[string]*Cell


type CellResponse struct {
	Value  string `json:"value"`
	Result string `json:"result"`
}

type SheetResponse map[string]CellResponse

type SetCell struct {
	Value string `json:"value"`
}
