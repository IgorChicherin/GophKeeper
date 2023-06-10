package models

import "time"

type CreateNoteRequest struct {
	Data     string `json:"data"`
	Metadata string `json:"metadata"`
	DataType string `json:"data_type"`
}

type DecodedNoteResponse struct {
	ID        int       `json:"id"`
	UserID    int       `json:"userId"`
	Data      string    `json:"data"`
	Metadata  string    `json:"metadata"`
	DataType  string    `json:"data_type"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
}
