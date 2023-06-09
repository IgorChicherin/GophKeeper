package models

type CreateNoteRequest struct {
	Data     string `json:"data"`
	Metadata string `json:"metadata"`
	DataType string `json:"data_type"`
}
